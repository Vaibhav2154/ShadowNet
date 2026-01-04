'use client'

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import {
  Server,
  Globe,
  Layers,
  LucideIcon
} from 'lucide-react'

// --- 1. Intro Section ---
export function IntroSection() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">What is a P2P Mesh VPN?</h1>
        <p className="text-lg text-neutral-400 leading-relaxed">
          ShadowNet represents a shift from traditional centralized VPN architectures to a decentralized,
          peer-to-peer (P2P) mesh topology. This section contrasts these approaches at a low technical level.
        </p>
      </div>

      <div className="grid md:grid-cols-2 gap-6">
        <Card className="border-neutral-800 bg-neutral-900/20">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Server className="w-5 h-5 text-neutral-500" />
              Hub-and-Spoke (Traditional)
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-neutral-300">
              In protocol terms, traffic from Node A to Node B is encapsulated and sent to a central Concentrator (C).
              C decapsulates, inspects routing tables, re-encapsulates, and forwards to B.
            </p>
            <div className="bg-black/50 p-4 rounded-lg font-mono text-xs text-neutral-400 border border-neutral-800">
              Latency(A→B) = Latency(A→C) + Latency(C→B) + Processing(C)
            </div>
            <ul className="space-y-2 text-sm text-neutral-400 list-disc list-inside">
              <li><strong>Bandwidth Bottleneck:</strong> Hub limits total network throughput.</li>
              <li><strong>Single Point of Failure:</strong> If Hub dies, network halts.</li>
              <li><strong>Hair-pinning:</strong> Traffic U-turns at the hub, wasting bandwidth.</li>
            </ul>
          </CardContent>
        </Card>

        <Card className="border-green-900/30 bg-green-950/10">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-white">
              <Globe className="w-5 h-5 text-green-500" />
              Full Mesh (ShadowNet)
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-neutral-300">
              Nodes maintain a routing table `O(n)` where `n` is peer count. Traffic is sent directly to the destination&apos;s
              discovered public endpoint using UDP.
            </p>
            <div className="bg-black/50 p-4 rounded-lg font-mono text-xs text-green-400 border border-green-900/30">
              Latency(A→B) = Latency(Direct Path)
            </div>
            <ul className="space-y-2 text-sm text-neutral-400 list-disc list-inside">
              <li><strong>Linear Scaling:</strong> Capacity grows with node count.</li>
              <li><strong>Resilience:</strong> Control plane failure affects only *new* connections.</li>
              <li><strong>Data Sovereignty:</strong> No intermediate server touches data.</li>
            </ul>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Architecture Deep Dive</CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="space-y-2">
            <h3 className="text-lg font-semibold text-white">The O(n²) Connection Problem</h3>
            <p className="text-neutral-400">
              In a full mesh of `N` nodes, there are potentially `N*(N-1)/2` connections. Maintaining state for thousands
              of tunnels is resource-intensive. ShadowNet uses <strong>Lazy Loading</strong>: tunnels (WireGuard sessions)
              are often initialized only when traffic is requested, though currently, we eagerly formulate the mesh for
              small networks.
            </p>
          </div>

          <div className="grid md:grid-cols-3 gap-4">
            <div className="p-4 bg-neutral-900/50 border border-neutral-800 rounded-lg">
              <div className="font-semibold text-white mb-1">Zero Trust</div>
              <p className="text-xs text-neutral-400">Identity is cryptographic (Public Key). IP addresses are ephemeral transport details.</p>
            </div>
            <div className="p-4 bg-neutral-900/50 border border-neutral-800 rounded-lg">
              <div className="font-semibold text-white mb-1">Self-Healing</div>
              <p className="text-xs text-neutral-400">If a direct path fails, STUN/Discovery logic re-runs to find new endpoints.</p>
            </div>
            <div className="p-4 bg-neutral-900/50 border border-neutral-800 rounded-lg">
              <div className="font-semibold text-white mb-1">NAT Agnostic</div>
              <p className="text-xs text-neutral-400">Works behind residential routers, CGNAT (Carrier Grade NAT), and LTE mobile networks.</p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

// --- 2. WireGuard Section ---
export function WireGuardSection() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">WireGuard Protocol Internals</h1>
        <p className="text-lg text-neutral-400 leading-relaxed">
          ShadowNet relies on the WireGuard protocol&apos;s rigorous formal verification and modern cryptography.
          It implements the <strong>Noise Protocol Framework</strong>, specifically a variant of <code>Noise_IKpsk2_25519_ChaChaPoly_BLAKE2s</code>.
        </p>
      </div>

      <div className="grid md:grid-cols-2 gap-6">
        <Card className="md:col-span-2">
          <CardHeader>
            <CardTitle>The Cryptographic Primitive Suite</CardTitle>
          </CardHeader>
          <CardContent className="grid md:grid-cols-4 gap-4">
            <div className="space-y-2">
              <div className="p-2 bg-blue-500/10 text-blue-400 rounded w-fit text-xs font-mono font-bold">Curve25519</div>
              <p className="text-xs text-neutral-400">ECDH (Elliptic Curve Diffie-Hellman) for key exchange. 32-byte public keys. Extremely fast scalar multiplication.</p>
            </div>
            <div className="space-y-2">
              <div className="p-2 bg-purple-500/10 text-purple-400 rounded w-fit text-xs font-mono font-bold">ChaCha20</div>
              <p className="text-xs text-neutral-400">Symmetric stream cipher. Encrypts the packet payload. chosen for speed on non-AES-NI (mobile) processors.</p>
            </div>
            <div className="space-y-2">
              <div className="p-2 bg-yellow-500/10 text-yellow-400 rounded w-fit text-xs font-mono font-bold">Poly1305</div>
              <p className="text-xs text-neutral-400">Message Authentication Code (MAC). Ensures packet integrity and authenticity.</p>
            </div>
            <div className="space-y-2">
              <div className="p-2 bg-red-500/10 text-red-400 rounded w-fit text-xs font-mono font-bold">BLAKE2s</div>
              <p className="text-xs text-neutral-400">Fast cryptographic hashing. Used for key derivation.</p>
            </div>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>The 1-RTT Handshake</CardTitle>
          <CardDescription>WireGuard establishes a secure session in just one round trip.</CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="space-y-4">
            <div className="relative pl-8 border-l-2 border-neutral-800 space-y-4">
              <div>
                <span className="absolute -left-2.5 top-0 w-5 h-5 bg-neutral-800 rounded-full flex items-center justify-center text-[10px] text-white">1</span>
                <h4 className="text-white font-medium text-sm">Initiation (A → B)</h4>
                <p className="text-xs text-neutral-400 mt-1">
                  Node A sends a handshake initiation. Contains ephemeral public key `E_a`, static public key `S_a` (encrypted),
                  and a MAC. Calculated as `Hash(ChainKey, ...)`.
                </p>
              </div>
              <div>
                <span className="absolute -left-2.5 top-0 w-5 h-5 bg-neutral-800 rounded-full flex items-center justify-center text-[10px] text-white mt-[4.5rem]">2</span>
                <h4 className="text-white font-medium text-sm">Response (B → A)</h4>
                <p className="text-xs text-neutral-400 mt-1">
                  Node B verifies MAC. If valid, computes shared secret. Sends back ephemeral `E_b` and auth tag.
                  Cookie is optionally sent back if B is under load (Active Defense).
                </p>
              </div>
              <div>
                <span className="absolute -left-2.5 top-0 w-5 h-5 bg-green-500 rounded-full flex items-center justify-center text-[10px] text-black font-bold mt-[9rem]">3</span>
                <h4 className="text-white font-medium text-sm">Transport Data (A → B)</h4>
                <p className="text-xs text-neutral-400 mt-1">
                  Keys are rotated every few minutes. Handshake confirms identity and establishes unique session keys
                  for sending and receiving.
                </p>
              </div>
            </div>
          </div>

          <div className="bg-neutral-900/50 p-4 rounded-lg border border-neutral-800">
            <h4 className="text-white font-medium mb-2 text-sm">Key Rotation & PFS</h4>
            <p className="text-xs text-neutral-300">
              WireGuard provides <strong>Perfect Forward Secrecy (PFS)</strong>. Session keys are ephemeral and strictly strictly timed.
              Even if an attacker steals your long-term private key later, they cannot decrypt past captured traffic because
              the ephemeral keys used to derive the session key are gone.
            </p>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

// --- 3. NAT Section ---
export function NATSection() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">NAT Traversal & UDP Hole Punching</h1>
        <p className="text-lg text-neutral-400 leading-relaxed">
          The &quot;Magic&quot; of ShadowNet is allowing two nodes behind different firewalls to connect directly.
          This is achieved via <strong>UDP Hole Punching</strong>, exploiting how stateful firewalls handle UDP.
        </p>
      </div>

      <div className="grid lg:grid-cols-2 gap-6">
        <Card className="h-full">
          <CardHeader>
            <CardTitle>The State Machine</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <ol className="space-y-4 text-sm text-neutral-300 list-decimal list-inside">
              <li className="pl-2">
                <strong>Binding Request:</strong> Node A sends UDP to STUN server.
                {/* Fixed: escaped < and > */}
                <br /><span className="text-xs text-neutral-500 ml-5 block">Firewall A creates entry: `PrivateIP:Port &lt;-&gt; PublicIP:Port`</span>
              </li>
              <li className="pl-2">
                <strong>Exchange:</strong> A and B swap their public endpoints via Control Plane (Side-channel).
              </li>
              <li className="pl-2">
                <strong>Simultaneous Open:</strong>
                <br /><span className="text-xs text-neutral-500 ml-5 block">A sends UDP to B. Firewall A opens gate. B&apos;s firewall drops it.</span>
                <span className="text-xs text-neutral-500 ml-5 block">B sends UDP to A. Firewall B opens gate.</span>
              </li>
              <li className="pl-2">
                <strong>Connection:</strong> A&apos;s packet arrives at B&apos;s firewall. Firewall B sees existing rule for &quot;sending to A&quot;, allows packet in.
              </li>
            </ol>
          </CardContent>
        </Card>

        <Card className="h-full">
          <CardHeader>
            <CardTitle>Keep-alives</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-neutral-300">
              NAT mappings expire quickly (often 60s for UDP). ShadowNet sends <strong>Persistent Keepalives</strong>.
            </p>
            <div className="space-y-2">
              <div className="flex justify-between text-xs text-neutral-400">
                <span>Default Keepalive</span>
                <span className="text-white font-mono">25 seconds</span>
              </div>
              <div className="w-full bg-neutral-800 h-1.5 rounded-full overflow-hidden">
                <div className="bg-blue-500 h-full w-[40%] animate-pulse"></div>
              </div>
              <p className="text-xs text-neutral-500">
                Ensures the router doesn&apos;t close the port mapping. If the mapping closes, the peer becomes unreachable
                until a new handshake (and perhaps a new STUN lookup) occurs.
              </p>
            </div>
          </CardContent>
        </Card>
      </div>

      <Card className="border-red-900/20 bg-red-950/5">
        <CardHeader>
          <CardTitle className="text-red-400">The Edge Case: Symmetric NAT</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-neutral-300">
            Some enterprise routers use <strong>Symmetric NAT</strong>. They assign a <em>different</em> public port
            for every destination IP.
          </p>
          <div className="mt-4 p-3 bg-black/40 rounded border border-red-900/30">
            <code className="text-xs font-mono text-red-300">
              A -&gt; STUN Server: 203.0.113.1:<strong>4001</strong><br />
              A -&gt; Node B     : 203.0.113.1:<strong>4005</strong> (Randomized!)
            </code>
          </div>
          <p className="text-xs text-neutral-400 mt-2">
            Hole punching fails because Node B sends traffic to port 4001, but Firewall A expects traffic on port 4005.
            Current workaround: <strong>Avoid Symmetric NATs</strong> or use a relay (DERP/TURN) in future versions.
          </p>
        </CardContent>
      </Card>
    </div>
  )
}

// --- 4. Control vs Data Plane ---
export function ControlDataSection() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Control Plane vs Data Plane</h1>
        <p className="text-lg text-neutral-400 leading-relaxed">
          ShadowNet employs a <strong>Split-Brain Architecture</strong>. The Control Plane is the &quot;Signal&quot;,
          the Mesh is the &quot;Media&quot;. This separation guarantees privacy and resilience.
        </p>
      </div>

      <div className="grid gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Control Plane (Signal)</CardTitle>
            <CardDescription>HTTP/JSON REST API - No persistent connections</CardDescription>
          </CardHeader>
          <CardContent className="grid md:grid-cols-2 gap-8">
            <div className="space-y-4">
              <h4 className="text-sm font-semibold text-white">Responsibilities</h4>
              <ul className="text-xs text-neutral-300 space-y-2">
                <li className="flex items-center gap-2"><div className="w-1.5 h-1.5 rounded-full bg-blue-500" /> Authenticating nodes via generic Auth Middleware.</li>
                <li className="flex items-center gap-2"><div className="w-1.5 h-1.5 rounded-full bg-blue-500" /> Storing Peer Map `(PubKey -&gt; [VirtualIP, PublicEndpoint])`.</li>
                <li className="flex items-center gap-2"><div className="w-1.5 h-1.5 rounded-full bg-blue-500" /> Distributing peer lists to authenticated nodes.</li>
              </ul>
            </div>
            <div className="bg-neutral-950 p-4 rounded-lg border border-neutral-800">
              <pre className="text-[10px] text-neutral-400 font-mono">
                {`POST /register HTTP/1.1
{
  "public_key": "base64...",
  "endpoint": "203.0.113.5:51820",
  "virtual_ip": "10.10.0.5"
}
-> 200 OK { peers: [...] }`}
              </pre>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Data Plane (Media)</CardTitle>
            <CardDescription>UDP / WireGuard - Purely Peer-to-Peer</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-neutral-300">
              Once Node A gets the peer list, it knows Node B&apos;s public endpoint. It programs the kernel (or userspace netstack)
              to route `10.10.0.B` packets into the encryptor.
            </p>
            <div className="p-4 bg-green-950/10 border border-green-900/30 rounded-lg">
              <h4 className="text-green-400 text-sm font-semibold mb-2">Failure Scenario: Control Plane Down</h4>
              <p className="text-xs text-neutral-400">
                If the Control Plane server crashes or is DDoS&apos;d:
                <br />
                1. <strong>Existing Connections:</strong> Unaffected. A and B continue to talk via UDP.
                <br />
                2. <strong>Roaming:</strong> If A changes IP, it cannot tell CP. B will lose connection to A eventually.
                <br />
                3. <strong>New Nodes:</strong> Cannot join.
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

// --- 5. Topology ---
export function TopologySection() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Network Topology & Routing</h1>
        <p className="text-lg text-neutral-400 leading-relaxed">
          How bits move from your application, through the virtual interface, across the internet, and into the destination application.
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>The Packet Lifecycle</CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="relative border-l border-neutral-800 ml-4 pl-8 space-y-8">
            <div className="relative">
              <span className="absolute -left-[37px] w-4 h-4 rounded-full bg-white ring-4 ring-black"></span>
              <h4 className="text-white font-mono text-sm">Application Layer</h4>
              <p className="text-xs text-neutral-400 mt-1">`curl 10.10.0.2:80` -&gt; Writes data to socket.</p>
            </div>
            <div className="relative">
              <span className="absolute -left-[37px] w-4 h-4 rounded-full bg-neutral-700 ring-4 ring-black"></span>
              <h4 className="text-white font-mono text-sm">OS Routing Table</h4>
              <p className="text-xs text-neutral-400 mt-1">
                Kernel sees destination `10.10.0.0/24`. Routing table says: `dev tun0`.
              </p>
            </div>
            <div className="relative">
              <span className="absolute -left-[37px] w-4 h-4 rounded-full bg-blue-600 ring-4 ring-black"></span>
              <h4 className="text-white font-mono text-sm">Virtual Interface (tun0)</h4>
              <p className="text-xs text-neutral-400 mt-1">
                ShadowNet (userspace) reads raw IP packet from `tun0` file descriptor.
              </p>
            </div>
            <div className="relative">
              <span className="absolute -left-[37px] w-4 h-4 rounded-full bg-purple-600 ring-4 ring-black"></span>
              <h4 className="text-white font-mono text-sm">WireGuard Encryptor</h4>
              <p className="text-xs text-neutral-400 mt-1">
                Look up Peer Public Key for IP `10.10.0.2`. Encrypt payload key ChaCha20. Add Poly1305 MAC.
              </p>
            </div>
            <div className="relative">
              <span className="absolute -left-[37px] w-4 h-4 rounded-full bg-green-600 ring-4 ring-black"></span>
              <h4 className="text-white font-mono text-sm">Physical Transport</h4>
              <p className="text-xs text-neutral-400 mt-1">
                Wrap encrypted blob in UDP. Send to `203.0.113.5:51820` (Node B&apos;s physical IP).
              </p>
            </div>
          </div>
        </CardContent>
      </Card>

      <div className="grid md:grid-cols-2 gap-4">
        <Card>
          <CardHeader><CardTitle className="text-sm">Metric Calculation</CardTitle></CardHeader>
          <CardContent>
            <p className="text-xs text-neutral-400">
              Currently, ShadowNet uses simple <strong>Static Routing</strong>. All peers are 1 hop away.
              In future versions with relaying, Dijkstra&apos;s algorithm would be used to calculate lowest-latency paths via relays.
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader><CardTitle className="text-sm">MTU & Fragmentation</CardTitle></CardHeader>
          <CardContent>
            <p className="text-xs text-neutral-400">
              WireGuard adds overhead (headers + auth tag).
              ShadowNet sets `tun0` MTU to <strong>1420 bytes</strong> (standard 1500 - 80 overhead) to avoid IP fragmentation, which degrades performance.
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

// --- 6. Security ---
export function SecuritySection() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Encryption & Security Model</h1>
        <p className="text-lg text-neutral-400 leading-relaxed">
          ShadowNet assumes the network is hostile. All traffic is authenticated and encrypted.
          We trust the math, not the infrastructure.
        </p>
      </div>

      <div className="grid md:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Threat Model</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Badge variant="destructive" className="mb-1">Passive Attacker</Badge>
              <p className="text-xs text-neutral-300">
                ISP or Wi-Fi snooper capturing packets.
                <br /><span className="text-neutral-500">Mitigation:</span> ChaCha20 Encryption. They see only random noise.
              </p>
            </div>
            <div className="space-y-2">
              <Badge variant="destructive" className="mb-1">Active Attacker</Badge>
              <p className="text-xs text-neutral-300">
                Man-in-the-Middle trying to modify packets.
                <br /><span className="text-neutral-500">Mitigation:</span> Poly1305 Auth Tag. Modified packets are dropped immediately.
              </p>
            </div>
            <div className="space-y-2">
              <Badge variant="destructive" className="mb-1">Replay Attacker</Badge>
              <p className="text-xs text-neutral-300">
                Resending valid old packets.
                <br /><span className="text-neutral-500">Mitigation:</span> Counter-based nonces. Old counters are rejected.
              </p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Identity & Trust</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-neutral-300">
              In ShadowNet, your <strong>Public Key IS your IP Address</strong> (conceptually).
              There are no passwords.
            </p>
            <div className="bg-neutral-900 p-4 rounded border border-neutral-800">
              <h4 className="text-white text-xs font-bold mb-2">Cryptographic ID Binding</h4>
              <code className="text-[10px] text-green-400 font-mono block break-all">
                NodeID: hash(PublicKey)<br />
                ACL: Allow traffic FROM 0xABC...123
              </code>
            </div>
            <p className="text-xs text-neutral-400 mt-2">
              To spoof a node, an attacker must possess its private key.
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

// --- 7. TUN ---
export function TUNSection() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">TUN/TAP Devices & Kernel Integration</h1>
        <p className="text-lg text-neutral-400 leading-relaxed">
          Understanding the boundary between Kernel Space and User Space.
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>TUN (Network Routing)</CardTitle>
          <CardDescription>Layer 3 Virtualization</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <p className="text-sm text-neutral-300">
            ShadowNet opens a file descriptor `/dev/net/tun`.
            Writing bytes to this FD transmits a packet to the OS network stack.
            Reading bytes from this FD receives packets routed by the OS to this interface.
          </p>
        </CardContent>
      </Card>

      <div className="grid md:grid-cols-2 gap-6">
        <Card className="bg-neutral-900/40">
          <CardHeader><CardTitle className="text-sm">Context Switching Cost</CardTitle></CardHeader>
          <CardContent>
            <p className="text-xs text-neutral-400">
              Every packet traversing ShadowNet moves: <br />
              `App (User) -&gt; Kernel -&gt; ShadowNet (User) -&gt; Kernel -&gt; NIC`.
              <br />
              This incurs overhead compared to kernel-mode WireGuard. Optimized by batch processing packets.
            </p>
          </CardContent>
        </Card>
        <Card className="bg-neutral-900/40">
          <CardHeader><CardTitle className="text-sm">Capabilities</CardTitle></CardHeader>
          <CardContent>
            <p className="text-xs text-neutral-400">
              Creating a TUN device requires `CAP_NET_ADMIN` on Linux. This is why the node binary usually needs `sudo`
              or explicit capabilities set via `setcap`.
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

// --- 8. Discovery ---
export function DiscoverySection() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-white mb-4">Peer Discovery Mechanics</h1>
        <p className="text-lg text-neutral-400 leading-relaxed">
          How disconnected nodes find each other on the vast internet.
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>The Registration Handshake</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-6">
            <div className="flex items-start gap-4">
              <div className="w-8 h-8 rounded bg-blue-900/30 flex items-center justify-center text-blue-400 font-mono text-sm shrink-0">1</div>
              <div>
                <h4 className="text-white text-sm font-semibold">Self-Discovery (STUN)</h4>
                <p className="text-xs text-neutral-400 mt-1">
                  Node sends binding request to Google/Mozilla STUN servers. Learn public IP.
                </p>
              </div>
            </div>
            <div className="flex items-start gap-4">
              <div className="w-8 h-8 rounded bg-blue-900/30 flex items-center justify-center text-blue-400 font-mono text-sm shrink-0">2</div>
              <div>
                <h4 className="text-white text-sm font-semibold">Announcement (HTTP)</h4>
                <p className="text-xs text-neutral-400 mt-1">
                  Node POSTs to Control Plane: &quot;I am Key X, available at IP Y&quot;.
                </p>
              </div>
            </div>
            <div className="flex items-start gap-4">
              <div className="w-8 h-8 rounded bg-blue-900/30 flex items-center justify-center text-blue-400 font-mono text-sm shrink-0">3</div>
              <div>
                <h4 className="text-white text-sm font-semibold">Convergence</h4>
                <p className="text-xs text-neutral-400 mt-1">
                  Control Plane responds with list of other active nodes.
                  Node immediately fires &quot;Hole Punch&quot; UDP packets to all of them.
                </p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

// --- 9. Integration ---
export function IntegrationSection() {
  return (
    <div className="space-y-8">
      <div>
        <div className="flex items-center gap-3 mb-4">
          <Layers className="w-10 h-10 text-white" />
          <h1 className="text-4xl font-bold text-white">How ShadowNet Integrates It All</h1>
        </div>
        <p className="text-lg text-neutral-400 leading-relaxed">
          We&apos;ve covered the primitives: WireGuard, STUN, TUN devices, NAT.
          Here is how ShadowNet orchestrates them into a cohesive system.
        </p>
      </div>

      <Card className="border-neutral-700 bg-neutral-900/50">
        <CardHeader>
          <CardTitle>Architecture Diagram</CardTitle>
        </CardHeader>
        <CardContent className="flex justify-center py-6">
          <div className="relative w-full max-w-2xl bg-black/50 p-6 rounded-lg border border-neutral-800 font-mono text-xs">
            <div className="flex justify-between items-center mb-12">
              <div className="border border-white/20 p-4 rounded text-center w-32">
                <div className="text-blue-400 mb-2">Control Plane</div>
                HTTP API
              </div>
              <div className="text-neutral-500 text-center">
                (Signaling Only)
                <div className="h-px bg-neutral-700 w-32 mx-auto my-2"></div>
                JSON / REST
              </div>
            </div>

            <div className="flex justify-between gap-12">
              <div className="border border-green-500/30 p-4 rounded flex-1">
                <div className="text-center font-bold text-white mb-4">Node A (You)</div>
                <div className="space-y-2">
                  <div className="bg-neutral-800 p-2 rounded text-neutral-300">Application</div>
                  <div className="text-center text-neutral-500">↓ write()</div>
                  <div className="bg-blue-900/30 border border-blue-500/30 p-2 rounded text-blue-200">TUN Interface</div>
                  <div className="text-center text-neutral-500">↓ read()</div>
                  <div className="bg-purple-900/30 border border-purple-500/30 p-2 rounded text-purple-200">ShadowNet Core</div>
                  <div className="text-center text-neutral-500">↓ encrypt()</div>
                  <div className="bg-yellow-900/30 border border-yellow-500/30 p-2 rounded text-yellow-200">UDP Socket</div>
                </div>
              </div>

              <div className="flex flex-col justify-end items-center pb-4 text-neutral-500 w-24">
                <div className="mb-2">Internet</div>
                <div className="w-full border-b border-dashed border-neutral-600"></div>
                <div className="my-2 text-[10px]">Encrypted UDP</div>
                <div className="w-full border-b border-dashed border-neutral-600"></div>
              </div>

              <div className="border border-green-500/30 p-4 rounded flex-1">
                <div className="text-center font-bold text-white mb-4">Node B (Peer)</div>
                <div className="space-y-2">
                  <div className="bg-neutral-800 p-2 rounded text-neutral-300">Application</div>
                  <div className="text-center text-neutral-500">↑ read()</div>
                  <div className="bg-blue-900/30 border border-blue-500/30 p-2 rounded text-blue-200">TUN Interface</div>
                  <div className="text-center text-neutral-500">↑ write()</div>
                  <div className="bg-purple-900/30 border border-purple-500/30 p-2 rounded text-purple-200">ShadowNet Core</div>
                  <div className="text-center text-neutral-500">↑ decrypt()</div>
                  <div className="bg-yellow-900/30 border border-yellow-500/30 p-2 rounded text-yellow-200">UDP Socket</div>
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <div className="space-y-6">
        <h3 className="text-2xl font-bold text-white">The Application Lifecycle</h3>

        <div className="grid gap-4">
          <Card className="relative overflow-hidden">
            <div className="absolute left-0 top-0 bottom-0 w-1 bg-white"></div>
            <CardHeader className="pb-2"><CardTitle className="text-base">1. Bootstrap</CardTitle></CardHeader>
            <CardContent>
              <p className="text-sm text-neutral-400">
                Executable starts. Loads private key from disk.
                Calls `ioctl` to create `tun0`. Sets IP `10.10.0.1`.
                Initializes UDP socket on port `51820`.
              </p>
            </CardContent>
          </Card>

          <Card className="relative overflow-hidden">
            <div className="absolute left-0 top-0 bottom-0 w-1 bg-blue-500"></div>
            <CardHeader className="pb-2"><CardTitle className="text-base">2. Discovery Loop</CardTitle></CardHeader>
            <CardContent>
              <p className="text-sm text-neutral-400">
                <strong>Every 30s:</strong>
                <br />1. Send STUN byte to Google. Get Public IP.
                <br />2. Send POST to Control Plane with Public IP + PubKey.
                <br />3. Receive JSON list of peers `[{`pubKey, endpoint_ip, virtual_ip`}, ...]`.
              </p>
            </CardContent>
          </Card>

          <Card className="relative overflow-hidden">
            <div className="absolute left-0 top-0 bottom-0 w-1 bg-purple-500"></div>
            <CardHeader className="pb-2"><CardTitle className="text-base">3. Mesh Convergence</CardTitle></CardHeader>
            <CardContent>
              <p className="text-sm text-neutral-400">
                For each peer in list:
                <br />- Update internal Routing Table: `10.10.0.X -&gt; PeerKey`.
                <br />- Send 0-byte UDP &quot;Keepalive&quot; to punch NAT.
                <br />- If application sends data, initiate WireGuard Handshake (Noise_IK).
              </p>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
