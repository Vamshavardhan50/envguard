  // src/components/SlideContent.tsx
"use client";

import React, { useState } from "react";
import { motion } from "framer-motion";
import { 
  AlertTriangle, Shield, CheckCircle, Terminal, 
  Settings, RefreshCw, Copy, Check, ExternalLink, Cpu
} from "lucide-react";

interface SlideContentProps {
  currentSlide: number;
  onActionClick: (actionType: string) => void;
  slideProgress: number;
}

// Bounding box with dashed border and corner handles matching the video editor preview window style
function TransformBox({ children, colorClass = "border-cyan-400/80" }: { children: React.ReactNode, colorClass?: string }) {
  return (
    <div className={`relative p-5 border-2 border-dashed ${colorClass} rounded-sm w-full max-w-md shadow-2xl bg-zinc-950/80 backdrop-blur-md z-20`}>
      {/* Corner Crop Handles */}
      <div className="absolute -top-1.5 -left-1.5 w-3 h-3 bg-white border border-black z-20" />
      <div className="absolute -top-1.5 -right-1.5 w-3 h-3 bg-white border border-black z-20" />
      <div className="absolute -bottom-1.5 -left-1.5 w-3 h-3 bg-white border border-black z-20" />
      <div className="absolute -bottom-1.5 -right-1.5 w-3 h-3 bg-white border border-black z-20" />
      {children}
    </div>
  );
}

// Kinetic typography wrappers for punchy reveals
function KineticText({ children, delay = 0, textClass = "text-white" }: { children: React.ReactNode, delay?: number, textClass?: string }) {
  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.95, y: 8, filter: "blur(4px)" }}
      animate={{ opacity: 1, scale: 1, y: 0, filter: "blur(0px)" }}
      transition={{ delay, duration: 0.35, ease: "easeOut" }}
      className={`font-black tracking-tight text-3xl md:text-5xl lowercase text-center leading-none ${textClass}`}
    >
      {children}
    </motion.div>
  );
}

export default function SlideContent({ currentSlide, onActionClick, slideProgress }: SlideContentProps) {
  const [copiedText, setCopiedText] = useState<string | null>(null);

  const copyToClipboard = (text: string, label: string) => {
    navigator.clipboard.writeText(text);
    setCopiedText(label);
    onActionClick("click");
    setTimeout(() => setCopiedText(null), 2000);
  };

  // Slide 4 Typewriter state
  const getInstallTerminalState = (progress: number) => {
    const cmd = "npm install -g envguard-bin";
    const chars = Math.min(cmd.length, Math.floor(Math.max(0, progress - 0.2) * 14));
    const typedCmd = cmd.substring(0, chars);
    
    let output = "";
    if (progress > 1.8) output += "added 108 packages, and audited 109 packages in 1.8s\n";
    if (progress > 2.4) output += "1 package is looking for funding. Run 'npm fund' for details\n\n";
    if (progress > 3.0) output += "✔ envguard successfully installed! Run 'envguard --help' to verify.\n";
    
    return { cmd: typedCmd, output };
  };

  // Slide 5 Typewriter state
  const getAuditTerminalState = (progress: number) => {
    const cmd = "envguard audit";
    const chars = Math.min(cmd.length, Math.floor(Math.max(0, progress - 0.2) * 12));
    const typedCmd = cmd.substring(0, chars);
    
    let output = "";
    if (progress > 1.2) output += "Scanning project files for environment key references...\n";
    if (progress > 2.0) output += "Comparing references against local .env file...\n\n";
    if (progress > 2.8) output += "  DATABASE_URL  ➡  Valid (matches database url rule)\n";
    if (progress > 3.5) output += "  PORT          ➡  Valid (matches number format)\n";
    if (progress > 4.2) output += "  JWT_SECRET    ➡  ✖ Missing in .env configuration!\n";
    if (progress > 4.9) output += "  OLD_API_KEY   ➡  ⚠ Unused (present in .env, not in code)\n\n";
    if (progress > 5.5) output += "✖ envguard audit failed: 1 missing variable, 1 unused variable.\n";
    
    return { cmd: typedCmd, output };
  };

  const installTerm = getInstallTerminalState(slideProgress);
  const auditTerm = getAuditTerminalState(slideProgress);

  const prevCmdLen = React.useRef(0);
  React.useEffect(() => {
    const activeCmd = currentSlide === 4 ? installTerm.cmd : (currentSlide === 5 ? auditTerm.cmd : "");
    if (activeCmd.length > prevCmdLen.current) {
      onActionClick("typewriter");
    }
    prevCmdLen.current = activeCmd.length;
  }, [installTerm.cmd, auditTerm.cmd, currentSlide, onActionClick]);

  return (
    <div className="flex-1 flex flex-col justify-center items-center px-6 md:px-24 z-20 text-white max-w-4xl w-full text-center relative select-none">
      
      {/* --- Slide 0: Teaser Hook (Solid Yellow Widescreen with Grid & Dynamic Waveforms) --- */}
      {currentSlide === 0 && (
        <div className="absolute inset-0 bg-[#facc15] z-10 flex flex-col justify-center items-center text-black px-6">
          {/* Editor grid lines background */}
          <div className="absolute inset-0 opacity-15 bg-[linear-gradient(rgba(0,0,0,0.15)_1px,transparent_1px),linear-gradient(90deg,rgba(0,0,0,0.15)_1px,transparent_1px)] bg-[size:30px_30px]" />
          
          <div className="space-y-4 z-20">
            <KineticText textClass="text-black font-black">ever deployed</KineticText>
            <KineticText delay={0.25} textClass="text-black font-black">a production app?</KineticText>
          </div>

          {/* Animated waveform columns at the bottom */}
          <div className="absolute bottom-12 left-0 w-full h-16 overflow-hidden flex items-end justify-around gap-1.5 opacity-25 px-16 pointer-events-none z-10">
            {Array.from({ length: 42 }).map((_, idx) => (
              <div
                key={idx}
                className="w-1.5 bg-black rounded-t-xs animate-pulse"
                style={{
                  height: `${20 + Math.random() * 80}%`,
                  animationDelay: `${idx * 0.04}s`,
                  animationDuration: `${0.7 + Math.random() * 1.1}s`
                }}
              />
            ))}
          </div>
        </div>
      )}

      {/* --- Slide 1: The Crash Incident (Dark Background, Red Crop Alert Box) --- */}
      {currentSlide === 1 && (
        <div className="flex flex-col items-center gap-6 relative">
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-80 h-80 bg-red-600/10 rounded-full filter blur-3xl animate-ping -z-10" />

          <div className="space-y-3 mb-4">
            <KineticText textClass="text-white">...only to see</KineticText>
            <KineticText delay={0.2} textClass="text-red-500 font-black underline decoration-wavy decoration-2">it crash instantly?</KineticText>
          </div>
          
          <TransformBox colorClass="border-red-500/80">
            <div className="flex items-center gap-2 border-b border-red-950 pb-2 mb-3 text-red-500/60 font-mono text-xxs">
              <span className="w-2 h-2 rounded-full bg-red-500 animate-ping" />
              <span>slack // alert-bot</span>
            </div>
            <div className="text-left font-mono text-xs text-red-300">
              <p className="font-bold text-red-500">[CRITICAL_DOWNTIME]</p>
              <p className="mt-1">app.js:144 Throw Error: Connection string DATABASE_URL is undefined</p>
              <p className="text-zinc-500 mt-2 text-[10px]">ReferenceError: DATABASE_URL is not defined in process.env</p>
            </div>
          </TransformBox>
        </div>
      )}

      {/* --- Slide 2: Pivot (Solid Cyan Widescreen with Morphing Emerald Blobs) --- */}
      {currentSlide === 2 && (
        <div className="absolute inset-0 bg-[#06b6d4] z-10 flex flex-col justify-center items-center text-black px-6 overflow-hidden">
          {/* Editor grid lines background */}
          <div className="absolute inset-0 opacity-15 bg-[linear-gradient(rgba(0,0,0,0.15)_1px,transparent_1px),linear-gradient(90deg,rgba(0,0,0,0.15)_1px,transparent_1px)] bg-[size:30px_30px]" />
          
          {/* Pulsing colored mesh lights */}
          <div className="absolute top-1/4 left-1/4 w-80 h-80 bg-emerald-300/40 rounded-full filter blur-3xl animate-pulse" />
          <div className="absolute bottom-1/4 right-1/4 w-80 h-80 bg-cyan-200/40 rounded-full filter blur-3xl animate-pulse" />

          <div className="space-y-4 z-20">
            <KineticText textClass="text-black font-black">you can skip</KineticText>
            <KineticText delay={0.2} textClass="text-black font-black">all of that.</KineticText>
          </div>
        </div>
      )}

      {/* --- Slide 3: Introducing envguard (Text overlay on top of centered Neon Shield & Rings) --- */}
      {currentSlide === 3 && (
        <div className="flex flex-col items-center gap-4">
          <div className="space-y-4">
            <KineticText textClass="text-cyan-400 font-extrabold">introducing envguard</KineticText>
            <KineticText delay={0.3} textClass="text-zinc-400 text-lg md:text-2xl font-bold font-mono">
              the eslint for environment variables.
            </KineticText>
          </div>
        </div>
      )}

      {/* --- Slide 4: Installation Terminal (Dashed Transform Box) --- */}
      {currentSlide === 4 && (
        <div className="flex flex-col items-center gap-6">
          <div className="space-y-3 mb-4">
            <KineticText textClass="text-white">install in one line. npm or pip.</KineticText>
          </div>

          <TransformBox colorClass="border-indigo-400/80">
            {/* Terminal Top bar */}
            <div className="flex items-center justify-between border-b border-zinc-900 pb-2 mb-2 text-zinc-500 font-mono text-xxs">
              <div className="flex gap-1">
                <div className="w-2 h-2 rounded-full bg-zinc-700" />
                <div className="w-2 h-2 rounded-full bg-zinc-700" />
                <div className="w-2 h-2 rounded-full bg-zinc-700" />
              </div>
              <span>bash -- install</span>
            </div>
            {/* CLI Area */}
            <div className="text-left font-mono text-xs h-32">
              <div className="flex items-center gap-1">
                <span className="text-zinc-500">$</span>
                <span className="text-cyan-400 font-bold">{installTerm.cmd}</span>
                <span className="w-1.5 h-3.5 bg-cyan-400 animate-pulse" />
              </div>
              <pre className="mt-2 text-xxs text-zinc-400 leading-normal whitespace-pre-wrap font-mono">
                {installTerm.output}
              </pre>
            </div>
          </TransformBox>
        </div>
      )}

      {/* --- Slide 5: CLI Scan (Dashed Transform Box) --- */}
      {currentSlide === 5 && (
        <div className="flex flex-col items-center gap-6">
          <div className="space-y-3 mb-4">
            <KineticText textClass="text-white">scan code. audit keys. zero leaks.</KineticText>
          </div>

          <TransformBox colorClass="border-emerald-400/80">
            {/* Terminal Header */}
            <div className="flex items-center justify-between border-b border-zinc-900 pb-2 mb-2 text-zinc-500 font-mono text-xxs">
              <div className="flex gap-1">
                <div className="w-2 h-2 rounded-full bg-zinc-700" />
                <div className="w-2 h-2 rounded-full bg-zinc-700" />
                <div className="w-2 h-2 rounded-full bg-zinc-700" />
              </div>
              <span>bash -- audit</span>
            </div>
            {/* CLI Area */}
            <div className="text-left font-mono text-xs h-48 overflow-y-auto">
              <div className="flex items-center gap-1">
                <span className="text-zinc-500">$</span>
                <span className="text-cyan-400 font-bold">{auditTerm.cmd}</span>
                <span className="w-1.5 h-3.5 bg-cyan-400 animate-pulse" />
              </div>
              <pre className="mt-2 text-xxs text-zinc-400 leading-normal whitespace-pre-wrap font-mono">
                {auditTerm.output}
              </pre>
            </div>
          </TransformBox>
        </div>
      )}

      {/* --- Slide 6: Rules Check (Dashed Transform Box) --- */}
      {currentSlide === 6 && (
        <div className="flex flex-col items-center gap-6">
          <div className="space-y-3 mb-4">
            <KineticText textClass="text-white">validate schemas, urls, and port defaults</KineticText>
          </div>

          <TransformBox colorClass="border-green-400/80">
            <div className="flex items-center justify-between border-b border-zinc-900 pb-2 mb-3 text-zinc-500 font-mono text-xxs">
              <span>config check</span>
              <span className="text-emerald-400 font-bold font-mono">.envguard.yaml</span>
            </div>
            
            <div className="space-y-2 text-left font-mono text-xxs">
              <div className="flex items-center justify-between p-2 bg-zinc-900/40 rounded border border-zinc-900/60">
                <span>DATABASE_URL: url (postgresql://)</span>
                <span className={`font-bold font-mono ${slideProgress > 1.0 ? "text-emerald-400" : "text-zinc-700"}`}>PASSED</span>
              </div>
              
              <div className="flex items-center justify-between p-2 bg-zinc-900/40 rounded border border-zinc-900/60">
                <span>PORT: number (default: 3000)</span>
                <span className={`font-bold font-mono ${slideProgress > 2.0 ? "text-emerald-400" : "text-zinc-700"}`}>PASSED</span>
              </div>
              
              <div className="flex items-center justify-between p-2 bg-zinc-900/40 rounded border border-zinc-900/60">
                <span>NODE_ENV: enum [dev, prod, test]</span>
                <span className={`font-bold font-mono ${slideProgress > 3.0 ? "text-emerald-400" : "text-zinc-700"}`}>PASSED</span>
              </div>
            </div>
          </TransformBox>
        </div>
      )}

      {/* --- Slide 7: Outro CTA (Transparent bg showing Warp-Speed Tunnel particles) --- */}
      {currentSlide === 7 && (
        <div className="flex flex-col justify-center items-center text-white px-6 w-full">
          {/* Subtle neon light aura */}
          <div className="absolute top-1/3 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[400px] h-[400px] bg-pink-500/10 rounded-full filter blur-3xl animate-pulse -z-10" />

          <div className="space-y-3 mb-8 z-20">
            <KineticText textClass="text-white font-black">secure your build.</KineticText>
            <KineticText delay={0.2} textClass="text-white font-black">grab envguard.</KineticText>
          </div>

          {/* Copyable code inputs with transparent glass effects */}
          <div className="space-y-3 w-full max-w-sm mb-8 z-20">
            <div className="bg-black/40 border border-white/10 rounded-xl p-3 flex items-center justify-between text-left backdrop-blur-md">
              <div className="font-mono text-xxs text-purple-200">
                <span className="text-white/40 block font-bold">NPM</span>
                <p className="text-white mt-0.5 font-mono">npm install -g envguard-bin</p>
              </div>
              <button
                onClick={() => copyToClipboard("npm install -g envguard-bin", "npm")}
                className="p-2 bg-white/10 hover:bg-white/20 transition rounded-lg text-white cursor-pointer"
              >
                {copiedText === "npm" ? <Check className="w-3.5 h-3.5 text-emerald-300" /> : <Copy className="w-3.5 h-3.5" />}
              </button>
            </div>

            <div className="bg-black/40 border border-white/10 rounded-xl p-3 flex items-center justify-between text-left backdrop-blur-md">
              <div className="font-mono text-xxs text-purple-200">
                <span className="text-white/40 block font-bold">PIP</span>
                <p className="text-white mt-0.5 font-mono">pip install envguard-bin</p>
              </div>
              <button
                onClick={() => copyToClipboard("pip install envguard-bin", "pip")}
                className="p-2 bg-white/10 hover:bg-white/20 transition rounded-lg text-white cursor-pointer"
              >
                {copiedText === "pip" ? <Check className="w-3.5 h-3.5 text-emerald-300" /> : <Copy className="w-3.5 h-3.5" />}
              </button>
            </div>
          </div>

          {/* GitHub links */}
          <div className="z-20">
            <a
              href="https://github.com/Vamshavardhan50/envguard"
              target="_blank"
              rel="noreferrer"
              onClick={() => onActionClick("click")}
              className="px-5 py-3 bg-black/40 hover:bg-black/60 transition rounded-lg font-bold text-xs flex items-center gap-2 border border-white/15 hover:border-white/25 cursor-pointer backdrop-blur-sm"
            >
              <svg className="w-4 h-4 fill-current text-white" viewBox="0 0 24 24">
                <path d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12"/>
              </svg>
              <span>Star on GitHub</span>
              <ExternalLink className="w-3.5 h-3.5 text-zinc-300" />
            </a>
          </div>
        </div>
      )}

    </div>
  );
}
