// src/utils/SynthManager.ts

export class SynthManager {
  private ctx: AudioContext | null = null;
  private isMuted: boolean = false;
  
  // Audio nodes
  private masterVolume: GainNode | null = null;
  private analyser: AnalyserNode | null = null;
  
  // Sidechain volume control (for pumping synth pad and arps)
  private sidechainVolume: GainNode | null = null;
  
  // Polyphonic pad oscillators
  private padOscs: OscillatorNode[] = [];
  private padGains: GainNode[] = [];
  private padFilter: BiquadFilterNode | null = null;
  
  // Arpeggiator volume
  private arpVolume: GainNode | null = null;
  private delayNode: DelayNode | null = null;
  private delayFeedback: GainNode | null = null;
  
  // Sequencer loop
  private arpIntervalId: any = null;
  private currentStep = 0;
  private theme = 0;
  private speed = 1.0;

  constructor() {
    // Synth will be initialized on first user interaction to bypass autoplay policy
  }

  public init() {
    if (typeof window === 'undefined') return;
    if (this.ctx) return;

    try {
      const AudioCtx = window.AudioContext || (window as any).webkitAudioContext;
      this.ctx = new AudioCtx();
      
      // Master volume control
      this.masterVolume = this.ctx.createGain();
      this.masterVolume.gain.setValueAtTime(0.12, this.ctx.currentTime); // keep it soft but clear

      // Setup analyser for real-time visual decibel meters
      this.analyser = this.ctx.createAnalyser();
      this.analyser.fftSize = 32;

      this.masterVolume.connect(this.analyser);
      this.analyser.connect(this.ctx.destination);

      // Sidechain volume node (connected to master)
      this.sidechainVolume = this.ctx.createGain();
      this.sidechainVolume.gain.setValueAtTime(1.0, this.ctx.currentTime);
      this.sidechainVolume.connect(this.masterVolume);

      // Setup delay loop for arpeggiator
      this.delayNode = this.ctx.createDelay(1.0);
      this.delayNode.delayTime.setValueAtTime(0.3, this.ctx.currentTime);
      this.delayFeedback = this.ctx.createGain();
      this.delayFeedback.gain.setValueAtTime(0.35, this.ctx.currentTime);

      this.arpVolume = this.ctx.createGain();
      this.arpVolume.gain.setValueAtTime(0.2, this.ctx.currentTime);

      // Connect arps to sidechain (so they pump!) and to delay loop
      this.arpVolume.connect(this.sidechainVolume);
      this.arpVolume.connect(this.delayNode);
      this.delayNode.connect(this.delayFeedback);
      this.delayFeedback.connect(this.delayNode);
      this.delayNode.connect(this.sidechainVolume);

      // Setup Polyphonic Pad
      this.setupPad();

      // Start arpeggiator/drum sequencer loop
      this.startSequencer();
    } catch (e) {
      console.warn("Failed to initialize Web Audio API:", e);
    }
  }

  private setupPad() {
    if (!this.ctx || !this.sidechainVolume) return;

    this.padFilter = this.ctx.createBiquadFilter();
    this.padFilter.type = 'lowpass';
    this.padFilter.frequency.setValueAtTime(220, this.ctx.currentTime); // Low cutoff for dark rumble
    this.padFilter.Q.setValueAtTime(4.0, this.ctx.currentTime);
    this.padFilter.connect(this.sidechainVolume);

    // Create 3 detuned pad oscillators for smooth chord shapes
    const waveforms: OscillatorType[] = ['sawtooth', 'triangle', 'sawtooth'];
    
    // Detuned frequencies for a tense C minor/dissonant diminished feel:
    // C2 (65.41), C#2 (69.30 - detuned dissonance!), G2 (98.00)
    const initialChord = [65.41, 69.30, 98.00];

    for (let i = 0; i < 3; i++) {
      const osc = this.ctx.createOscillator();
      const gain = this.ctx.createGain();
      
      osc.type = waveforms[i];
      osc.frequency.setValueAtTime(initialChord[i], this.ctx.currentTime);
      
      // Detune slightly for chorused width
      osc.detune.setValueAtTime((i - 1) * 8, this.ctx.currentTime);

      gain.gain.setValueAtTime(0.05, this.ctx.currentTime);
      
      osc.connect(gain);
      gain.connect(this.padFilter);
      
      osc.start();
      this.padOscs.push(osc);
      this.padGains.push(gain);
    }
  }

  private setPadChord(frequencies: number[]) {
    if (!this.ctx || this.padOscs.length < 3) return;
    const now = this.ctx.currentTime;
    
    frequencies.forEach((freq, idx) => {
      const osc = this.padOscs[idx];
      if (osc) {
        // Glide frequencies smoothly to new notes (portamento)
        osc.frequency.exponentialRampToValueAtTime(freq, now + 1.2);
      }
    });
  }

  private startSequencer() {
    if (this.arpIntervalId) clearInterval(this.arpIntervalId);
    
    // Clock tick interval
    const interval = Math.max(50, Math.floor(150 / this.speed));
    this.arpIntervalId = setInterval(() => {
      this.tick();
    }, interval);
  }

  public setSpeed(speed: number) {
    this.speed = speed;
    if (this.ctx) {
      this.startSequencer();
    }
  }

  private tick() {
    if (!this.ctx || this.isMuted || this.ctx.state === 'suspended') return;
    
    const now = this.ctx.currentTime;

    // Slow sweep the lowpass filter dynamically to create a "breathing" filter sweep
    if (this.padFilter) {
      const breathingCutoff = 180 + Math.sin(now * 1.5) * 50;
      this.padFilter.frequency.setValueAtTime(breathingCutoff, now);
    }

    // Trigger arpeggiator plucks based on active theme
    this.playArpStep(now);

    // Trigger procedural drum heartbeat / ticking clock
    this.playDrumStep(now);

    this.currentStep = (this.currentStep + 1) % 16;
  }

  private playArpStep(now: number) {
    if (!this.ctx || !this.arpVolume) return;

    // Tense minor/diminished arpeggiation patterns to match suspense
    // Chords:
    // 0: Teaser -> atmosphere only
    // 1: The Crash -> high pitch alarms
    // 2: Solution Hook -> C minor plucks
    // 3: Installation -> D minor diminished plucks
    // 4: CLI Demo -> F minor diminished scan
    // 5: Health Schema -> G minor suspense
    // 6: Outro -> C minor resolved
    const chords: { [key: number]: number[] } = {
      0: [], // Silence (pure tension drone)
      1: [523.25, 554.37, 783.99, 880.00], // High alarms C5, C#5, G5, A5
      2: [130.81, 155.56, 185.00, 196.00], // C minor/diminished (C3, Eb3, F#3, G3)
      3: [146.83, 174.61, 207.65, 220.00], // D diminished (D3, F3, Ab3, A3)
      4: [174.61, 207.65, 246.94, 261.63], // F diminished (F3, Ab3, B3, C4)
      5: [196.00, 233.08, 277.18, 293.66], // G diminished (G3, Bb3, Db4, D4)
      6: [130.81, 155.56, 196.00, 261.63]  // Resolved C minor
    };

    const currentNotes = chords[this.theme] || chords[2];
    if (currentNotes.length === 0) return; // Silent arp on teaser

    // Syncopated tension arpeggiation masks
    const patterns: { [key: number]: number[] } = {
      1: [1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0], // Fast heartbeat alerts
      2: [1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 1, 0, 0, 1, 0], // Syncopated pluck
      3: [1, 0, 1, 0, 0, 1, 0, 1, 0, 1, 0, 0, 1, 0, 1, 0], // Tense build
      4: [1, 1, 0, 1, 0, 1, 1, 0, 1, 1, 0, 1, 0, 1, 1, 0], // Intricate scan
      5: [1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0], // Slower suspense
      6: [1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0]  // Ending pad pluck
    };

    const activePattern = patterns[this.theme] || patterns[2];
    const shouldPlay = activePattern[this.currentStep] === 1;

    if (shouldPlay) {
      const noteFreq = currentNotes[this.currentStep % currentNotes.length];
      
      const osc = this.ctx.createOscillator();
      const gain = this.ctx.createGain();
      const filter = this.ctx.createBiquadFilter();

      // Pluck Synth sound: triangle/saw with decaying filter sweep
      osc.type = this.theme === 1 ? 'sawtooth' : 'triangle';
      osc.frequency.setValueAtTime(noteFreq, now);

      filter.type = 'lowpass';
      filter.frequency.setValueAtTime(1800, now);
      filter.frequency.exponentialRampToValueAtTime(250, now + 0.15); // rapid pluck decay
      filter.Q.setValueAtTime(4.0, now);

      gain.gain.setValueAtTime(0, now);
      gain.gain.linearRampToValueAtTime(this.theme === 1 ? 0.04 : 0.07, now + 0.015);
      gain.gain.exponentialRampToValueAtTime(0.0001, now + 0.3);

      osc.connect(filter);
      filter.connect(gain);
      gain.connect(this.arpVolume);

      osc.start(now);
      osc.stop(now + 0.35);
    }
  }

  private playDrumStep(now: number) {
    // Ticking Clock Sound: Plays constant clock ticks on step 0, 4, 8, 12 to build urgency
    if (this.currentStep % 4 === 0) {
      this.playClockTick(now);
    }

    // Heartbeat double-thump Kick
    // plays: THUMP-THUMP (0, 1) and THUMP-THUMP (8, 9)
    if (this.theme > 0) {
      if (this.currentStep === 0 || this.currentStep === 8) {
        this.playKick(now, false); // Main Thump
      }
      if (this.currentStep === 1 || this.currentStep === 9) {
        this.playKick(now + 0.02, true); // Ghost Sub-Thump
      }
    }
    
    // Slats or glitches (Slide 1/The Crash has extra alerts)
    if (this.theme === 1 && this.currentStep === 6) {
      this.playSnare(now);
    }
  }

  private playClockTick(now: number) {
    if (!this.ctx || !this.masterVolume) return;
    const osc = this.ctx.createOscillator();
    const gain = this.ctx.createGain();
    
    // High-pitched clock foley tick
    osc.type = 'sine';
    osc.frequency.setValueAtTime(2800, now);
    
    gain.gain.setValueAtTime(0.012, now);
    gain.gain.exponentialRampToValueAtTime(0.0001, now + 0.01);
    
    osc.connect(gain);
    gain.connect(this.masterVolume);
    
    osc.start(now);
    osc.stop(now + 0.015);
  }

  private playKick(now: number, isSub: boolean = false) {
    if (!this.ctx || !this.masterVolume) return;
    const osc = this.ctx.createOscillator();
    const gain = this.ctx.createGain();
    
    osc.type = 'sine';
    // Lower pitch for sub thumps
    osc.frequency.setValueAtTime(isSub ? 80 : 110, now);
    osc.frequency.exponentialRampToValueAtTime(0.001, now + 0.1);
    
    gain.gain.setValueAtTime(isSub ? 0.12 : 0.18, now);
    gain.gain.exponentialRampToValueAtTime(0.0001, now + 0.1);
    
    osc.connect(gain);
    gain.connect(this.masterVolume);
    
    osc.start(now);
    osc.stop(now + 0.11);

    // Dynamic sidechain pumping on kick thumps
    if (this.sidechainVolume && !isSub) {
      this.sidechainVolume.gain.setValueAtTime(1.0, now);
      this.sidechainVolume.gain.linearRampToValueAtTime(0.4, now + 0.015); // Dip
      this.sidechainVolume.gain.exponentialRampToValueAtTime(1.0, now + 0.18); // Release
    }
  }

  private playSnare(now: number) {
    if (!this.ctx || !this.masterVolume) return;
    const osc = this.ctx.createOscillator();
    const gain = this.ctx.createGain();
    
    osc.type = 'triangle';
    osc.frequency.setValueAtTime(160, now);
    
    gain.gain.setValueAtTime(0.04, now);
    gain.gain.exponentialRampToValueAtTime(0.0001, now + 0.07);
    
    osc.connect(gain);
    gain.connect(this.masterVolume);
    
    osc.start(now);
    osc.stop(now + 0.08);
  }

  public setTheme(themeId: number) {
    this.theme = themeId;
    if (!this.ctx) return;
    
    const now = this.ctx.currentTime;
    
    // Chord mappings for polyphonic pad glide
    const chords: { [key: number]: number[] } = {
      0: [65.41, 69.30, 98.00], // C minor detuned rub (C2, C#2, G2)
      1: [65.41, 69.30, 92.50], // C diminished rub (C2, C#2, F#2)
      2: [65.41, 77.78, 98.00], // C minor (C2, Eb2, G2)
      3: [73.42, 87.31, 110.00], // D diminished (D2, F2, Ab2)
      4: [87.31, 103.83, 130.81], // F diminished (F2, Ab2, C3)
      5: [55.00, 65.41, 82.41],   // A minor (A1, C2, E2)
      6: [65.41, 77.78, 98.00]    // C minor resolve
    };
    
    const currentChord = chords[themeId] || chords[2];
    this.setPadChord(currentChord);
  }

  public playTypewriter() {
    if (!this.ctx || this.isMuted || this.ctx.state === 'suspended') return;
    
    const now = this.ctx.currentTime;
    const osc = this.ctx.createOscillator();
    const gain = this.ctx.createGain();
    
    osc.type = 'triangle';
    osc.frequency.setValueAtTime(1400, now);
    osc.frequency.exponentialRampToValueAtTime(320, now + 0.025);
    
    gain.gain.setValueAtTime(0, now);
    gain.gain.linearRampToValueAtTime(0.02, now + 0.002);
    gain.gain.exponentialRampToValueAtTime(0.0001, now + 0.025);
    
    osc.connect(gain);
    gain.connect(this.masterVolume!);
    
    osc.start(now);
    osc.stop(now + 0.03);
  }

  public playWhoosh() {
    if (!this.ctx || this.isMuted || this.ctx.state === 'suspended') return;
    
    const now = this.ctx.currentTime;
    const osc = this.ctx.createOscillator();
    const gain = this.ctx.createGain();
    const filter = this.ctx.createBiquadFilter();
    
    osc.type = 'sawtooth';
    osc.frequency.setValueAtTime(80, now);
    osc.frequency.exponentialRampToValueAtTime(750, now + 0.35);
    
    filter.type = 'lowpass';
    filter.frequency.setValueAtTime(120, now);
    filter.frequency.exponentialRampToValueAtTime(1600, now + 0.35);
    
    gain.gain.setValueAtTime(0, now);
    gain.gain.linearRampToValueAtTime(0.05, now + 0.08);
    gain.gain.exponentialRampToValueAtTime(0.0001, now + 0.4);
    
    osc.connect(filter);
    filter.connect(gain);
    gain.connect(this.masterVolume!);
    
    osc.start(now);
    osc.stop(now + 0.45);
  }

  public playClick() {
    if (!this.ctx || this.isMuted || this.ctx.state === 'suspended') return;

    const now = this.ctx.currentTime;
    const osc = this.ctx.createOscillator();
    const gain = this.ctx.createGain();

    osc.type = 'sine';
    osc.frequency.setValueAtTime(800, now);
    osc.frequency.exponentialRampToValueAtTime(1200, now + 0.05);

    gain.gain.setValueAtTime(0, now);
    gain.gain.linearRampToValueAtTime(0.05, now + 0.005);
    gain.gain.exponentialRampToValueAtTime(0.0001, now + 0.08);

    osc.connect(gain);
    gain.connect(this.masterVolume!);

    osc.start(now);
    osc.stop(now + 0.1);
  }

  public playGlitch() {
    if (!this.ctx || this.isMuted || this.ctx.state === 'suspended') return;
    
    const now = this.ctx.currentTime;
    const osc = this.ctx.createOscillator();
    const gain = this.ctx.createGain();
    
    osc.type = 'sawtooth';
    osc.frequency.setValueAtTime(120, now);
    osc.frequency.setValueAtTime(60, now + 0.04);
    osc.frequency.setValueAtTime(250, now + 0.08);
    osc.frequency.setValueAtTime(80, now + 0.12);
    
    gain.gain.setValueAtTime(0, now);
    gain.gain.linearRampToValueAtTime(0.08, now + 0.01);
    gain.gain.exponentialRampToValueAtTime(0.0001, now + 0.18);
    
    osc.connect(gain);
    gain.connect(this.masterVolume!);
    
    osc.start(now);
    osc.stop(now + 0.2);
  }

  public toggleMute(): boolean {
    if (!this.ctx || !this.masterVolume) return true;
    
    this.isMuted = !this.isMuted;
    const targetVal = this.isMuted ? 0 : 0.12;
    this.masterVolume.gain.setTargetAtTime(targetVal, this.ctx.currentTime, 0.05);
    
    return this.isMuted;
  }

  public getMutedState() {
    return this.isMuted;
  }

  public resume() {
    if (this.ctx && this.ctx.state === 'suspended') {
      this.ctx.resume();
    }
  }

  public pause() {
    // Sequencer remains active
  }

  public getVolumeLevel(): number {
    if (!this.analyser || this.isMuted) return 0;
    const array = new Uint8Array(this.analyser.frequencyBinCount);
    this.analyser.getByteFrequencyData(array);
    let sum = 0;
    for (let i = 0; i < array.length; i++) {
      sum += array[i];
    }
    return (sum / array.length) / 255;
  }

  public destroy() {
    if (this.arpIntervalId) clearInterval(this.arpIntervalId);
    this.padOscs.forEach(o => {
      try { o.stop(); } catch(e){}
    });
    this.padOscs = [];
    if (this.ctx) {
      this.ctx.close();
    }
  }
}
