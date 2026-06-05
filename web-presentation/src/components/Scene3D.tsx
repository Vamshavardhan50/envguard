// src/components/Scene3D.tsx
"use client";

import React, { useRef, useMemo } from "react";
import { Canvas, useFrame, useThree } from "@react-three/fiber";
import { Points, PointMaterial } from "@react-three/drei";
import * as THREE from "three";

// --- Sub-Component: Camera Controller Rig ---
interface CameraRigProps {
  theme: number;
}

function CameraRig({ theme }: CameraRigProps) {
  const { camera } = useThree();
  const currentLookAt = useRef(new THREE.Vector3(0, 0, 0));

  useFrame((state) => {
    const time = state.clock.getElapsedTime();
    const targetPos = new THREE.Vector3(0, 0, 5);
    const targetLookAt = new THREE.Vector3(0, 0, 0);

    // Slide 0: Teaser (Solid Yellow Screen)
    if (theme === 0) {
      targetPos.set(0, 0, 6.2);
      targetLookAt.set(0, 0, 0);
    }
    // Slide 1: The Crash (Handheld shake, red alerts)
    else if (theme === 1) {
      const shakeFreq = 16.0;
      const shakeAmp = 0.15;
      const shakeX = Math.sin(time * shakeFreq) * Math.cos(time * shakeFreq * 0.8) * shakeAmp;
      const shakeY = Math.cos(time * shakeFreq * 1.3) * Math.sin(time * shakeFreq * 0.9) * shakeAmp;
      const shakeZ = Math.sin(time * shakeFreq * 0.6) * shakeAmp * 0.5;

      targetPos.set(-1.6 + shakeX, 0.6 + shakeY, 3.8 + shakeZ);
      targetLookAt.set(0.3, 0.15, -0.4);
    }
    // Slide 2: Pivot (Solid Cyan Screen)
    else if (theme === 2) {
      targetPos.set(0, 0.1, 4.5);
      targetLookAt.set(0, 0.1, 0);
    }
    // Slide 3: Intro Shield (Shield Centered, expanding rings)
    else if (theme === 3) {
      targetPos.set(0, 0, 4.5);
      targetLookAt.set(0, 0, 0);
    }
    // Slide 4: Install (Shield Shrunk Left, terminal typing)
    else if (theme === 4) {
      targetPos.set(-2.2, 0.8, 3.8);
      targetLookAt.set(0.8, 0.2, -0.6);
    }
    // Slide 5: CLI Scan (Data streams, scanning)
    else if (theme === 5) {
      targetPos.set(-2.2, 0.8, 3.8);
      targetLookAt.set(0.8, 0.2, -0.6);
    }
    // Slide 6: Rules Check (Orbit security nodes)
    else if (theme === 6) {
      const radius = 4.2;
      const orbitSpeed = 0.15;
      const angle = time * orbitSpeed;
      targetPos.set(Math.sin(angle) * radius, 0.8, Math.cos(angle) * radius);
      targetLookAt.set(0, 0, 0);
    }
    // Slide 7: Outro CTA (Warp-Speed Star Tunnel flying past camera)
    else if (theme === 7) {
      targetPos.set(0, 0, 4.0);
      targetLookAt.set(0, 0, 0);
    }

    camera.position.lerp(targetPos, 0.08);
    currentLookAt.current.lerp(targetLookAt, 0.08);
    camera.lookAt(currentLookAt.current);
  });

  return null;
}

// --- Sub-Component 1: Glitch Matrix ---
interface GlitchMatrixProps {
  isAlert?: boolean;
}

function GlitchMatrix({ isAlert = false }: GlitchMatrixProps) {
  const pointsRef = useRef<THREE.Points>(null);
  
  const particleCount = 1500;
  const [positions, speeds] = useMemo(() => {
    const pos = new Float32Array(particleCount * 3);
    const spd = new Float32Array(particleCount);
    for (let i = 0; i < particleCount; i++) {
      pos[i * 3] = (Math.random() - 0.5) * 15; // X
      pos[i * 3 + 1] = (Math.random() - 0.5) * 15; // Y
      pos[i * 3 + 2] = (Math.random() - 0.5) * 15; // Z
      spd[i] = 0.05 + Math.random() * 0.15;
    }
    return [pos, spd];
  }, []);

  useFrame((state) => {
    if (!pointsRef.current) return;
    const time = state.clock.getElapsedTime();
    const posAttr = pointsRef.current.geometry.attributes.position;
    
    const speedMult = isAlert ? 2.5 : 0.6;
    const waveMult = isAlert ? 3.0 : 0.8;
    
    for (let i = 0; i < particleCount; i++) {
      let y = posAttr.getY(i) - speeds[i] * speedMult * 0.5;
      if (y < -7.5) y = 7.5;
      posAttr.setY(i, y);

      const x = posAttr.getX(i);
      const wave = Math.sin(y * 2 + time * (isAlert ? 25 : 10)) * 0.15 * Math.sin(time * 3) * waveMult;
      posAttr.setX(i, x + wave * 0.03);
    }
    posAttr.needsUpdate = true;
    pointsRef.current.rotation.y = time * 0.03;
  });

  return (
    <Points ref={pointsRef} positions={positions} stride={3} frustumCulled={false}>
      <PointMaterial
        transparent
        color={isAlert ? "#ef4444" : "#991b1b"}
        size={isAlert ? 0.08 : 0.05}
        sizeAttenuation={true}
        depthWrite={false}
        blending={THREE.AdditiveBlending}
      />
    </Points>
  );
}

// --- Sub-Component: Expansion Rings (Inside-Out Waves) ---
function ExpansionRings() {
  const count = 4;
  const ringsRef = useRef<THREE.Group>(null);

  useFrame((state) => {
    if (!ringsRef.current) return;
    const time = state.clock.getElapsedTime();
    const children = ringsRef.current.children;
    
    for (let i = 0; i < children.length; i++) {
      const mesh = children[i] as THREE.Mesh;
      const progress = ((time * 0.6 + i / count) % 1.0); // 0 to 1
      const scale = 0.1 + progress * 3.4;
      mesh.scale.set(scale, scale, scale);
      
      const mat = mesh.material as THREE.MeshBasicMaterial;
      if (mat) {
        mat.opacity = Math.max(0, (1 - progress) * 0.4);
      }
    }
  });

  return (
    <group ref={ringsRef}>
      {Array.from({ length: count }).map((_, i) => (
        <mesh key={i} rotation={[Math.PI / 2, 0, 0]}>
          <ringGeometry args={[0.95, 1.0, 64]} />
          <meshBasicMaterial
            color="#06b6d4"
            transparent
            opacity={0.3}
            side={THREE.DoubleSide}
            depthWrite={false}
          />
        </mesh>
      ))}
    </group>
  );
}

// --- Sub-Component 2: Neon Hologram Shield ---
interface ShieldProps {
  shrink: boolean;
}
function HologramShield({ shrink }: ShieldProps) {
  const groupRef = useRef<THREE.Group>(null);
  const ringRef = useRef<THREE.Mesh>(null);

  useFrame((state) => {
    if (!groupRef.current) return;
    const time = state.clock.getElapsedTime();
    
    groupRef.current.rotation.y = time * 0.4;
    groupRef.current.rotation.z = Math.sin(time * 0.5) * 0.1;
    
    if (ringRef.current) {
      ringRef.current.rotation.x = time * -0.6;
    }

    const targetScale = shrink ? 0.35 : 1.0;
    const targetY = shrink ? 1.8 : 0;
    const targetX = shrink ? -2.2 : 0;

    groupRef.current.scale.lerp(new THREE.Vector3(targetScale, targetScale, targetScale), 0.1);
    groupRef.current.position.lerp(new THREE.Vector3(targetX, targetY, 0), 0.1);
  });

  const shieldShape = useMemo(() => {
    const shape = new THREE.Shape();
    shape.moveTo(0, 1.2);
    shape.quadraticCurveTo(0.8, 1.1, 1.0, 0.4);
    shape.quadraticCurveTo(0.9, -0.6, 0, -1.4);
    shape.quadraticCurveTo(-0.9, -0.6, -1.0, 0.4);
    shape.quadraticCurveTo(-0.8, 1.1, 0, 1.2);
    return shape;
  }, []);

  const extrudeSettings = {
    steps: 1,
    depth: 0.15,
    bevelEnabled: true,
    bevelThickness: 0.05,
    bevelSize: 0.05,
    bevelSegments: 3
  };

  return (
    <group ref={groupRef}>
      <mesh ref={ringRef}>
        <torusGeometry args={[1.7, 0.02, 8, 64]} />
        <meshBasicMaterial color="#06b6d4" transparent opacity={0.3} wireframe />
      </mesh>

      {/* Inside-to-outside expansion rings */}
      {!shrink && <ExpansionRings />}
      
      <mesh castShadow receiveShadow>
        <extrudeGeometry args={[shieldShape, extrudeSettings]} />
        <meshStandardMaterial
          color="#0891b2"
          wireframe
          transparent
          opacity={0.7}
          emissive="#06b6d4"
          emissiveIntensity={1.5}
        />
      </mesh>

      <mesh scale={[0.9, 0.9, 0.9]}>
        <extrudeGeometry args={[shieldShape, { ...extrudeSettings, depth: 0.08 }]} />
        <meshBasicMaterial
          color="#0e7490"
          transparent
          opacity={0.15}
        />
      </mesh>
    </group>
  );
}

// --- Sub-Component 3: Code Scanning / Data Streams ---
function DataStreams() {
  const groupRef = useRef<THREE.Group>(null);
  const streamCount = 20;

  const streams = useMemo(() => {
    return Array.from({ length: streamCount }).map(() => ({
      y: (Math.random() - 0.5) * 6,
      z: (Math.random() - 0.5) * 4,
      xStart: 5 + Math.random() * 5,
      speed: 0.08 + Math.random() * 0.15,
      size: 0.05 + Math.random() * 0.15,
      color: Math.random() > 0.5 ? "#22c55e" : "#06b6d4"
    }));
  }, []);

  useFrame(() => {
    if (!groupRef.current) return;
    const children = groupRef.current.children;
    streams.forEach((stream, idx) => {
      const child = children[idx] as THREE.Mesh;
      if (!child) return;
      
      child.position.x -= stream.speed;
      if (child.position.x < -6) {
        child.position.x = stream.xStart;
      }
    });
  });

  return (
    <group ref={groupRef}>
      {streams.map((stream, idx) => (
        <mesh key={idx} position={[stream.xStart, stream.y, stream.z]}>
          <boxGeometry args={[stream.size * 2, stream.size, stream.size]} />
          <meshBasicMaterial color={stream.color} transparent opacity={0.6} />
        </mesh>
      ))}
    </group>
  );
}

// --- Sub-Component 4: Constellation Security Nodes ---
function SecurityNodes() {
  const groupRef = useRef<THREE.Group>(null);
  
  const nodes = useMemo(() => {
    return [
      { pos: [-2, 1.5, 0], label: "DATABASE_URL" },
      { pos: [2, 1, -1], label: "PORT" },
      { pos: [0, -1, 1], label: "NODE_ENV" },
      { pos: [-3, -1.2, -1], label: "STRIPE_KEY" },
      { pos: [3, -1.5, 0], label: "JWT_SECRET" }
    ];
  }, []);

  const connections = useMemo(() => {
    const list: THREE.Line[] = [];
    nodes.forEach((n1, i) => {
      nodes.slice(i + 1).forEach((n2, j) => {
        const p1 = new THREE.Vector3(...n1.pos);
        const p2 = new THREE.Vector3(...n2.pos);
        const distance = p1.distanceTo(p2);
        if (distance <= 4.5) {
          const lineGeom = new THREE.BufferGeometry().setFromPoints([p1, p2]);
          const lineMaterial = new THREE.LineBasicMaterial({
            color: "#10b981",
            transparent: true,
            opacity: 0.25
          });
          list.push(new THREE.Line(lineGeom, lineMaterial));
        }
      });
    });
    return list;
  }, [nodes]);

  useFrame((state) => {
    if (!groupRef.current) return;
    const time = state.clock.getElapsedTime();
    groupRef.current.rotation.y = Math.sin(time * 0.2) * 0.3;
    groupRef.current.rotation.x = Math.cos(time * 0.15) * 0.1;
  });

  return (
    <group ref={groupRef}>
      {nodes.map((node, idx) => (
        <group key={idx} position={node.pos as [number, number, number]}>
          <mesh>
            <sphereGeometry args={[0.22, 16, 16]} />
            <meshStandardMaterial
              color="#22c55e"
              emissive="#22c55e"
              emissiveIntensity={0.8}
              roughness={0.1}
            />
          </mesh>
          <mesh rotation={[Math.PI / 2, 0, 0]}>
            <torusGeometry args={[0.35, 0.01, 8, 32]} />
            <meshBasicMaterial color="#4ade80" transparent opacity={0.4} />
          </mesh>
        </group>
      ))}
      {connections.map((lineObj, idx) => (
        <primitive key={idx} object={lineObj} />
      ))}
    </group>
  );
}

// --- Sub-Component 5: Warp-Speed Star Tunnel (Inside-Out Flying) ---
function WarpTunnel() {
  const pointsRef = useRef<THREE.Points>(null);
  const particleCount = 1200;

  const [positions, speeds] = useMemo(() => {
    const pos = new Float32Array(particleCount * 3);
    const spd = new Float32Array(particleCount);
    for (let i = 0; i < particleCount; i++) {
      const angle = Math.random() * Math.PI * 2;
      const radius = 1.0 + Math.random() * 3.0; // Keep tunnel center hollow
      pos[i * 3] = Math.cos(angle) * radius; // X
      pos[i * 3 + 1] = Math.sin(angle) * radius; // Y
      pos[i * 3 + 2] = -20 + Math.random() * 25; // Z (spread out from -20 to 5)
      spd[i] = 0.15 + Math.random() * 0.25; // Flying speed
    }
    return [pos, spd];
  }, []);

  useFrame((state) => {
    if (!pointsRef.current) return;
    const time = state.clock.getElapsedTime();
    const posAttr = pointsRef.current.geometry.attributes.position;

    for (let i = 0; i < particleCount; i++) {
      let z = posAttr.getZ(i) + speeds[i] * 1.5;
      if (z > 6) {
        z = -20; // reset to the end
      }
      posAttr.setZ(i, z);

      // Rotate particles slightly for a twisting tunnel visual
      const x = posAttr.getX(i);
      const y = posAttr.getY(i);
      const rotSpeed = 0.002;
      const cosR = Math.cos(rotSpeed);
      const sinR = Math.sin(rotSpeed);
      posAttr.setX(i, x * cosR - y * sinR);
      posAttr.setY(i, x * sinR + y * cosR);
    }
    posAttr.needsUpdate = true;
    pointsRef.current.rotation.z = time * 0.05;
  });

  return (
    <Points ref={pointsRef} positions={positions} stride={3} frustumCulled={false}>
      <PointMaterial
        transparent
        color="#d8b4fe" // Glowing violet stars
        size={0.07}
        sizeAttenuation={true}
        depthWrite={false}
        blending={THREE.AdditiveBlending}
      />
    </Points>
  );
}

// --- Main Scene Container with Camera Control ---
interface SceneProps {
  theme: number;
}
export default function Scene3D({ theme }: SceneProps) {
  // Determine if Three.js elements should render or be hidden behind solid colored pages
  // Slides 0 (yellow) and 2 (cyan) are full solid colors, so we skip Three.js render overlays for performance.
  const isSolidScreen = theme === 0 || theme === 2;
  const bgClass = theme === 7 ? "bg-[#180a2b]" : "bg-black"; // Rich dark purple for outro CTA

  return (
    <div className={`absolute inset-0 w-full h-full -z-10 ${bgClass} overflow-hidden transition-colors duration-500`}>
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_center,rgba(9,9,11,0.2)_0%,rgba(0,0,0,1)_85%)] pointer-events-none z-10" />

      {!isSolidScreen && (
        <Canvas
          camera={{ position: [0, 0, 5], fov: 60 }}
          dpr={[1, 1.5]}
        >
          <ambientLight intensity={0.5} />
          <pointLight position={[10, 10, 10]} intensity={1.5} />
          
          {theme === 1 && <GlitchMatrix isAlert={true} />}
          
          {theme === 3 && <HologramShield shrink={false} />}
          
          {theme === 4 && (
            <HologramShield shrink={true} />
          )}

          {theme === 5 && (
            <>
              <HologramShield shrink={true} />
              <DataStreams />
            </>
          )}
          
          {theme === 6 && <SecurityNodes />}

          {theme === 7 && <WarpTunnel />}

          <CameraRig theme={theme} />
        </Canvas>
      )}
    </div>
  );
}
