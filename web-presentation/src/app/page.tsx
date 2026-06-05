// src/app/page.tsx
"use client";

import React, { useState, useEffect, useRef } from "react";
import Scene3D from "@/components/Scene3D";
import SlideContent from "@/components/SlideContent";
import { SynthManager } from "@/utils/SynthManager";
import { 
  Volume2, VolumeX, Play, Pause, RotateCcw, Sparkles, 
  Terminal, SkipForward, SkipBack, Film, Radio
} from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";

const slides = [
  { id: 0, name: "Teaser", start: 0, end: 4, color: "border-yellow-900/30 bg-yellow-950/20" },
  { id: 1, name: "The Crash", start: 4, end: 10, color: "border-red-900/30 bg-red-950/10" },
  { id: 2, name: "Pivot", start: 10, end: 14, color: "border-cyan-900/30 bg-cyan-950/10" },
  { id: 3, name: "Intro Shield", start: 14, end: 18, color: "border-blue-900/30 bg-blue-950/10" },
  { id: 4, name: "Install", start: 18, end: 22, color: "border-indigo-900/30 bg-indigo-950/10" },
  { id: 5, name: "CLI Scan", start: 22, end: 28, color: "border-emerald-900/30 bg-emerald-950/10" },
  { id: 6, name: "Rules Check", start: 28, end: 34, color: "border-green-900/30 bg-green-950/10" },
  { id: 7, name: "Outro CTA", start: 34, end: 45, color: "border-purple-900/30 bg-purple-950/10" }
];
const totalDuration = 45;

export default function Home() {
  const [hasStarted, setHasStarted] = useState(false);
  const [currentTime, setCurrentTime] = useState(0);
  const [audioEnabled, setAudioEnabled] = useState(false);
  const [isPlaying, setIsPlaying] = useState(false);
  const [playbackSpeed, setPlaybackSpeed] = useState(1);
  const [decibels, setDecibels] = useState<number[]>([0, 0]);
  const [isBlinking, setIsBlinking] = useState(false);
  const [flashColor, setFlashColor] = useState<string | null>(null);
  const [flashKey, setFlashKey] = useState(0);
  
  const synthRef = useRef<SynthManager | null>(null);
  const timelineRef = useRef<HTMLDivElement>(null);

  // Initialize SynthManager on Mount
  useEffect(() => {
    synthRef.current = new SynthManager();
    return () => {
      if (synthRef.current) {
        synthRef.current.destroy();
      }
    };
  }, []);

  // Flashing REC Dot Timer
  useEffect(() => {
    const interval = setInterval(() => {
      setIsBlinking((prev) => !prev);
    }, 1000);
    return () => clearInterval(interval);
  }, []);

  // Frame accurate timeline counter using requestAnimationFrame
  useEffect(() => {
    if (!isPlaying || !hasStarted) return;
    
    let lastTime = performance.now();
    let frameId: number;

    const update = (now: number) => {
      const delta = (now - lastTime) / 1000 * playbackSpeed;
      lastTime = now;

      setCurrentTime((prev) => {
        const next = prev + delta;
        if (next >= totalDuration) {
          return 0; // Loop back
        }
        return next;
      });
      
      frameId = requestAnimationFrame(update);
    };

    frameId = requestAnimationFrame(update);
    return () => cancelAnimationFrame(frameId);
  }, [isPlaying, hasStarted, playbackSpeed]);

  // Read volume level from SynthManager analyser
  useEffect(() => {
    if (!hasStarted) return;
    let frameId: number;

    const poll = () => {
      if (synthRef.current) {
        const volume = synthRef.current.getVolumeLevel();
        // Introduce small random stereo drift for aesthetic L/R differences
        const left = Math.min(1.0, volume * 1.6);
        const right = Math.min(1.0, volume * (1.3 + Math.random() * 0.3));
        setDecibels([left, right]);
      }
      frameId = requestAnimationFrame(poll);
    };

    frameId = requestAnimationFrame(poll);
    return () => cancelAnimationFrame(frameId);
  }, [hasStarted]);

  // Auto detect active slide based on current timeline index
  const activeSlide = slides.find(s => currentTime >= s.start && currentTime < s.end)?.id ?? 7;

  // Sync Audio theme and trigger transition screen flash when activeSlide changes
  useEffect(() => {
    if (!hasStarted) return;

    // 1. Trigger transition color flash
    const flashColors: { [key: number]: string } = {
      0: "rgba(250, 204, 21, 0.75)",  // Yellow
      1: "rgba(239, 68, 68, 0.85)",   // Red
      2: "rgba(6, 182, 212, 0.8)",    // Cyan
      3: "rgba(59, 130, 246, 0.75)",   // Blue
      4: "rgba(99, 102, 241, 0.75)",   // Indigo
      5: "rgba(16, 185, 129, 0.75)",   // Emerald
      6: "rgba(34, 197, 94, 0.75)",    // Green
      7: "rgba(168, 85, 247, 0.8)"    // Purple
    };

    const color = flashColors[activeSlide] || "rgba(255, 255, 255, 0.8)";
    setFlashColor(color);
    setFlashKey(prev => prev + 1);

    // 2. Synchronize synth themes
    if (synthRef.current) {
      const themeMap: { [key: number]: number } = {
        0: 0, // Teaser -> detuned rub C minor
        1: 1, // Crash -> detuned rub C diminished
        2: 2, // Pivot -> C minor pluck
        3: 3, // Intro Shield -> D diminished pluck
        4: 4, // Install -> F diminished pluck
        5: 5, // CLI Scan -> A minor scan
        6: 6, // Rules Check -> C minor resolve
        7: 6  // Outro CTA -> C minor resolve
      };
      synthRef.current.setTheme(themeMap[activeSlide] ?? 6);
      
      // Slide 1 has a glitch sound, others get a whoosh sound
      if (activeSlide === 1) {
        synthRef.current.playGlitch();
      } else {
        synthRef.current.playWhoosh();
      }
    }
  }, [activeSlide, hasStarted]);

  // Synchronize synth arpeggiator speed on playback speed changes
  useEffect(() => {
    if (synthRef.current && hasStarted) {
      synthRef.current.setSpeed(playbackSpeed);
    }
  }, [playbackSpeed, hasStarted]);

  // Synchronize music play/pause state with timeline playback
  useEffect(() => {
    if (!synthRef.current || !hasStarted) return;
    if (isPlaying) {
      synthRef.current.resume();
    } else {
      synthRef.current.pause();
    }
  }, [isPlaying, hasStarted]);

  const startPresentation = () => {
    setHasStarted(true);
    setAudioEnabled(true);
    setIsPlaying(true);

    if (synthRef.current) {
      synthRef.current.init();
      synthRef.current.resume();
      synthRef.current.setTheme(0);
      synthRef.current.setSpeed(playbackSpeed);
    }
  };

  const toggleMute = () => {
    if (synthRef.current) {
      const isMuted = synthRef.current.toggleMute();
      setAudioEnabled(!isMuted);
    }
  };

  const handleActionClick = (actionType: string) => {
    if (!synthRef.current) return;
    if (actionType === "click") {
      synthRef.current.playClick();
    } else if (actionType === "glitch") {
      synthRef.current.playGlitch();
    } else if (actionType === "typewriter") {
      synthRef.current.playTypewriter();
    }
  };

  const handleTimelineClick = (e: React.MouseEvent<HTMLDivElement>) => {
    if (!timelineRef.current) return;
    const rect = timelineRef.current.getBoundingClientRect();
    const clickX = e.clientX - rect.left;
    const pct = Math.max(0, Math.min(1, clickX / rect.width));
    
    setCurrentTime(pct * totalDuration);
    handleActionClick("click");
  };

  const skipForward = () => {
    const nextIdx = (activeSlide + 1) % slides.length;
    setCurrentTime(slides[nextIdx].start);
    handleActionClick("click");
  };

  const skipBack = () => {
    const prevIdx = (activeSlide - 1 + slides.length) % slides.length;
    setCurrentTime(slides[prevIdx].start);
    handleActionClick("click");
  };

  // Convert time to standard 30fps Timecode (HH:MM:SS:FF)
  const getTimecode = (timeInSecs: number) => {
    const totalFrames = Math.floor(timeInSecs * 30);
    const frames = totalFrames % 30;
    const totalSeconds = Math.floor(timeInSecs);
    const seconds = totalSeconds % 60;
    const minutes = Math.floor(totalSeconds / 60) % 60;
    const hours = Math.floor(totalSeconds / 3600);
    
    return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}:${String(frames).padStart(2, '0')}`;
  };

  return (
    <main className="relative flex flex-col min-h-screen w-full select-none overflow-hidden bg-black font-sans text-white">
      
      {/* 3D WebGL Canvas Layer */}
      <Scene3D theme={activeSlide} />

      {/* Editor Grid Lines Overlay for dark slides (visible behind text but on top of Canvas) */}
      <div className="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.012)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.012)_1px,transparent_1px)] bg-[size:40px_40px] pointer-events-none z-10" />

      {/* --- Overlay 1: Onboarding Gate --- */}
      <AnimatePresence>
        {!hasStarted && (
          <motion.div 
            initial={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="absolute inset-0 bg-black/95 backdrop-blur-2xl z-50 flex flex-col justify-center items-center text-center p-6"
          >
            <motion.div 
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              transition={{ delay: 0.1, duration: 0.6 }}
              className="max-w-lg bg-zinc-950/70 border border-zinc-800 p-10 rounded-2xl shadow-2xl backdrop-blur-lg relative overflow-hidden"
            >
              <div className="absolute top-0 left-0 w-full h-[2px] bg-gradient-to-r from-red-500 via-cyan-500 to-purple-500" />
              
              <div className="flex justify-center mb-6">
                <div className="p-4 bg-cyan-500/10 rounded-full border border-cyan-500/20 text-cyan-400">
                  <Film className="w-10 h-10 animate-pulse" />
                </div>
              </div>
              
              <h1 className="text-4xl font-black text-white mb-4 tracking-tight leading-none uppercase">
                envguard <span className="text-cyan-400">Cinematic</span>
              </h1>
              
              <p className="text-sm text-zinc-400 mb-8 leading-relaxed">
                Experience the 3D product demonstration sequence with synchronized soundscapes, automated command simulations, and real-time audio analysis.
              </p>

              <button
                onClick={startPresentation}
                className="w-full py-4 bg-gradient-to-r from-cyan-500 via-emerald-500 to-purple-500 hover:from-cyan-600 hover:to-purple-600 text-black font-extrabold rounded-xl shadow-lg transition active:scale-98 tracking-wider uppercase"
              >
                Initiate Sequence & Start Audio
              </button>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* --- Cinematic 2.39:1 Aspect Ratio Bars (Video Editor Vibe) --- */}
      {hasStarted && (
        <>
          <div className="absolute top-0 left-0 w-full h-8 bg-black/90 border-b border-zinc-900/60 z-40 flex items-center justify-between px-6 text-xxs font-mono tracking-widest text-zinc-500 pointer-events-none">
            <span>MONITOR: 4K UHD 2.39:1 CROP</span>
            <span>COLOR SPACE: S-LOG3 NEON</span>
            <span>LUT: ENVGUARD_PREMIUM.CUBE</span>
          </div>
          <div className="absolute bottom-0 left-0 w-full h-8 bg-black/90 border-t border-zinc-900/60 z-40 flex items-center justify-between px-6 text-xxs font-mono tracking-widest text-zinc-500 pointer-events-none">
            <span>FPS: 30.00 (CONSTANT)</span>
            <span>AUDIO SOURCE: WEB_SYNTH_BUS</span>
            <span>RENDER ENGINE: WEBGL_THREE</span>
          </div>
        </>
      )}

      {/* --- Viewfinder Indicators (Top Corners) --- */}
      {hasStarted && (
        <div className="absolute top-12 left-6 right-6 flex justify-between items-start z-30 pointer-events-none">
          {/* Top Left: REC Status */}
          <div className="flex flex-col gap-1.5 bg-black/60 backdrop-blur-md px-4 py-2.5 rounded-lg border border-zinc-900 font-mono text-xs">
            <div className="flex items-center gap-2">
              <span className={`w-2.5 h-2.5 rounded-full bg-red-600 ${isBlinking && isPlaying ? "opacity-100" : "opacity-30"}`} />
              <span className="font-bold text-white tracking-widest uppercase">REC</span>
            </div>
            <div className="text-zinc-400 text-xxs tracking-wider uppercase flex items-center gap-1.5">
              <Radio className="w-3 h-3 text-cyan-400 animate-pulse" />
              <span>LIVE DEMO SEQUENCE</span>
            </div>
          </div>

          {/* Top Right: Timecode Counter */}
          <div className="flex flex-col items-end gap-1.5 bg-black/60 backdrop-blur-md px-4 py-2.5 rounded-lg border border-zinc-900 font-mono text-xs">
            <div className="text-emerald-400 font-bold tracking-wider">
              {getTimecode(currentTime)}
            </div>
            <div className="text-zinc-500 text-xxs tracking-wider uppercase">
              TOTAL: {getTimecode(totalDuration)}
            </div>
          </div>
        </div>
      )}

      {/* --- Top Navigation Header --- */}
      {hasStarted && (
        <header className="px-6 pt-12 pb-4 flex items-center justify-between border-b border-zinc-900/40 bg-black/20 backdrop-blur-sm z-30">
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 bg-cyan-500/10 border border-cyan-500/20 rounded-lg flex items-center justify-center text-cyan-400 shadow-inner">
              <Terminal className="w-4 h-4" />
            </div>
            <div>
              <span className="font-black text-sm tracking-widest uppercase text-white">envguard</span>
              <span className="text-xxs text-zinc-500 block font-mono">v1.0.3 cinematic renderer</span>
            </div>
          </div>

          {/* Volume Mute Toggle */}
          <button
            onClick={toggleMute}
            className={`p-2.5 rounded-lg border transition flex items-center gap-2 text-xs font-semibold cursor-pointer z-40 ${audioEnabled ? "bg-cyan-500/10 border-cyan-500/20 text-cyan-400" : "bg-zinc-900/50 border-zinc-800/80 text-zinc-400 hover:text-white"}`}
          >
            {audioEnabled ? (
              <>
                <Volume2 className="w-4 h-4 animate-bounce" />
                <span className="hidden sm:inline">AUDIO CONNECTED</span>
              </>
            ) : (
              <>
                <VolumeX className="w-4 h-4" />
                <span className="hidden sm:inline">AUDIO MUTED</span>
              </>
            )}
          </button>
        </header>
      )}

      {/* --- Main Interactive Content Slide Overlay --- */}
      {hasStarted && (
        <div className="flex-1 flex flex-col justify-center items-center relative py-12">
          {/* Audio Decibel Meters (Floating Left) */}
          <div className="absolute left-6 bottom-36 hidden md:flex flex-row gap-2.5 items-end bg-black/60 border border-zinc-900 p-4 rounded-xl backdrop-blur-md z-30 w-16 h-48 pointer-events-none">
            {decibels.map((channelVal, chIdx) => (
              <div key={chIdx} className="flex-1 flex flex-col justify-end h-full gap-[2px]">
                {Array.from({ length: 12 }).map((_, stepIdx) => {
                  const stepLimit = (12 - stepIdx) / 12;
                  const isActive = channelVal >= stepLimit;
                  
                  // Color code: green at bottom, yellow middle, red top
                  let colorClass = "bg-zinc-900";
                  if (isActive) {
                    if (stepIdx < 3) colorClass = "bg-red-500 shadow-[0_0_8px_rgba(239,68,68,0.5)]";
                    else if (stepIdx < 6) colorClass = "bg-yellow-500 shadow-[0_0_8px_rgba(234,179,8,0.5)]";
                    else colorClass = "bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.5)]";
                  }
                  
                  return (
                    <div 
                      key={stepIdx} 
                      className={`w-full h-2 rounded-xs transition-colors duration-75 ${colorClass}`}
                    />
                  );
                })}
                <span className="text-zinc-600 font-mono text-[8px] text-center mt-1">
                  {chIdx === 0 ? "L" : "R"}
                </span>
              </div>
            ))}
          </div>

          <SlideContent 
            currentSlide={activeSlide} 
            onActionClick={handleActionClick} 
            slideProgress={currentTime - slides[activeSlide].start}
          />
        </div>
      )}

      {/* --- Bottom Cinematic Timeline & Scrubber Panel --- */}
      {hasStarted && (
        <footer className="px-6 pb-12 pt-4 bg-zinc-950/80 border-t border-zinc-900 backdrop-blur-md flex flex-col gap-4 z-30">
          
          {/* Interactive Multi-track Timeline Scrubber */}
          <div className="flex flex-col gap-1.5 w-full">
            <div className="flex items-center justify-between text-[10px] font-mono text-zinc-500">
              <span className="flex items-center gap-1.5">
                <Film className="w-3.5 h-3.5 text-cyan-400" />
                TIMELINE TRACKS
              </span>
              <span>SCENE: {slides[activeSlide].name}</span>
            </div>

            {/* Scrubber Container */}
            <div 
              ref={timelineRef}
              onClick={handleTimelineClick}
              className="relative h-6 bg-zinc-900/60 border border-zinc-800/80 rounded-lg cursor-ew-resize overflow-hidden flex"
            >
              {/* Timeline segment tracks */}
              {slides.map((s) => {
                const widthPct = ((s.end - s.start) / totalDuration) * 100;
                const isActive = activeSlide === s.id;
                return (
                  <div
                    key={s.id}
                    style={{ width: `${widthPct}%` }}
                    className={`h-full border-r border-zinc-800/60 relative flex flex-col justify-between p-1 select-none transition ${s.color} ${isActive ? "opacity-100" : "opacity-60"}`}
                  >
                    <span className="text-[9px] font-mono font-bold tracking-tight text-zinc-400 truncate">
                      {s.name}
                    </span>
                    <div className={`h-[2px] w-full rounded ${isActive ? "bg-cyan-500" : "bg-transparent"}`} />
                  </div>
                );
              })}

              {/* Playback seeker line (The cursor head) */}
              <div 
                style={{ left: `${(currentTime / totalDuration) * 100}%` }}
                className="absolute top-0 bottom-0 w-[2px] bg-red-500 pointer-events-none shadow-[0_0_8px_rgba(239,68,68,1)] z-10"
              >
                <div className="absolute top-0 left-1/2 -translate-x-1/2 w-3 h-3 bg-red-500 rotate-45 border border-white" />
              </div>
            </div>
          </div>

          {/* Media Control Toolbar */}
          <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
            
            {/* Speed Controller & Status */}
            <div className="flex items-center gap-4">
              <div className="flex items-center gap-1.5">
                <button
                  onClick={() => {
                    setIsPlaying(!isPlaying);
                    handleActionClick("click");
                  }}
                  className={`p-2.5 rounded-lg transition-all cursor-pointer ${isPlaying ? "bg-cyan-500 text-black hover:bg-cyan-400" : "bg-zinc-900 text-zinc-400 hover:text-white border border-zinc-800"}`}
                  title={isPlaying ? "Pause Sequence" : "Play Sequence"}
                >
                  {isPlaying ? <Pause className="w-4 h-4 fill-current" /> : <Play className="w-4 h-4 fill-current" />}
                </button>
                <button
                  onClick={skipBack}
                  className="p-2.5 bg-zinc-900 hover:bg-zinc-800 border border-zinc-800 transition rounded-lg text-zinc-400 hover:text-white cursor-pointer"
                  title="Previous Scene"
                >
                  <SkipBack className="w-4 h-4" />
                </button>
                <button
                  onClick={skipForward}
                  className="p-2.5 bg-zinc-900 hover:bg-zinc-800 border border-zinc-800 transition rounded-lg text-zinc-400 hover:text-white cursor-pointer"
                  title="Next Scene"
                >
                  <SkipForward className="w-4 h-4" />
                </button>
                <button
                  onClick={() => {
                    setCurrentTime(0);
                    handleActionClick("click");
                  }}
                  className="p-2.5 bg-zinc-900 hover:bg-zinc-800 border border-zinc-800 transition rounded-lg text-zinc-400 hover:text-white cursor-pointer"
                  title="Rewind to start"
                >
                  <RotateCcw className="w-4 h-4" />
                </button>
              </div>

              <div className="h-6 w-[1px] bg-zinc-900 hidden sm:block" />

              <div className="hidden sm:flex items-center gap-1.5">
                <span className="text-[10px] text-zinc-500 font-mono">SPEED:</span>
                {[1, 1.5, 2].map((speed) => (
                  <button
                    key={speed}
                    onClick={() => {
                      setPlaybackSpeed(speed);
                      handleActionClick("click");
                    }}
                    className={`px-2 py-0.5 text-xxs font-mono rounded cursor-pointer ${playbackSpeed === speed ? "bg-cyan-500/10 text-cyan-400 border border-cyan-500/20" : "text-zinc-500 hover:text-zinc-300"}`}
                  >
                    {speed}x
                  </button>
                ))}
              </div>
            </div>

            {/* Interactive Stage Progression Bar & Description */}
            <div className="flex items-center gap-2">
              <span className="text-xxs text-zinc-500 font-mono uppercase tracking-widest">
                {isPlaying ? "AUTOMATED PREVIEW IN PROGRESS" : "SEQUENCE PAUSED"}
              </span>
              <span className="w-2 h-2 rounded-full bg-cyan-400 animate-ping" />
            </div>

          </div>
        </footer>
      )}
      {/* Flash Cut Overlay */}
      <AnimatePresence>
        {flashColor && (
          <motion.div
            key={flashKey}
            initial={{ opacity: 1 }}
            animate={{ opacity: 0 }}
            transition={{ duration: 0.25, ease: "easeOut" }}
            style={{ backgroundColor: flashColor }}
            className="absolute inset-0 pointer-events-none z-50"
          />
        )}
      </AnimatePresence>
    </main>
  );
}
