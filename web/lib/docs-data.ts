import {
  Network,
  Shield,
  Wifi,
  Lock,
  Server,
  Globe,
  ArrowLeftRight,
  Zap,
  Layers,
  LucideIcon
} from 'lucide-react'
import {
  IntroSection,
  WireGuardSection,
  NATSection,
  ControlDataSection,
  TopologySection,
  SecuritySection,
  TUNSection,
  DiscoverySection,
  IntegrationSection
} from '@/components/docs/sections'

export type SectionId =
  | 'intro'
  | 'wireguard'
  | 'nat'
  | 'control-data'
  | 'topology'
  | 'security'
  | 'tun'
  | 'discovery'
  | 'integration'

export interface DocSection {
  id: SectionId
  title: string
  icon: LucideIcon
  component: React.ComponentType
}

export const DOCS_SECTIONS: DocSection[] = [
  { id: 'intro', title: 'What is P2P Mesh VPN?', icon: Network, component: IntroSection },
  { id: 'wireguard', title: 'WireGuard Protocol', icon: Shield, component: WireGuardSection },
  { id: 'nat', title: 'NAT Traversal', icon: Wifi, component: NATSection },
  { id: 'control-data', title: 'Control vs Data Plane', icon: Server, component: ControlDataSection },
  { id: 'topology', title: 'Network Topology', icon: Globe, component: TopologySection },
  { id: 'security', title: 'Encryption & Security', icon: Lock, component: SecuritySection },
  { id: 'tun', title: 'TUN/TAP Devices', icon: ArrowLeftRight, component: TUNSection },
  { id: 'discovery', title: 'Peer Discovery', icon: Zap, component: DiscoverySection },
  { id: 'integration', title: 'ShadowNet Integration', icon: Layers, component: IntegrationSection },
]
