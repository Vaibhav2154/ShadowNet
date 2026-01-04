'use client'

import { useEffect, useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { NetworkTopology } from '@/components/network-topology'
import {
  Activity,
  Users,
  Server,
  Clock,
  Wifi,
  WifiOff,
  RefreshCw,
  Network,
  Zap,
  Globe,
  Home,
  Settings,
  BarChart3,
  Shield
} from 'lucide-react'
import { timeAgo } from '@/lib/utils'

interface PeerInfo {
  id: string
  wg_public_key: string
  endpoint_ip: string
  endpoint_port: number
  last_seen: string
}

interface Metrics {
  total_peers: number
  active_peers: number
  uptime: string
  timestamp: string
}

const CONTROL_PLANE_URL = process.env.NEXT_PUBLIC_CONTROLPLANE_URL || 'http://localhost:8080'

export default function Dashboard() {
  const [peers, setPeers] = useState<PeerInfo[]>([])
  const [metrics, setMetrics] = useState<Metrics | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [lastUpdate, setLastUpdate] = useState<Date>(new Date())

  const fetchData = async () => {
    try {
      setError(null)

      const [peersRes, metricsRes] = await Promise.all([
        fetch(`${CONTROL_PLANE_URL}/peers`),
        fetch(`${CONTROL_PLANE_URL}/metrics`)
      ])

      if (!peersRes.ok) throw new Error('Failed to fetch peers')
      if (!metricsRes.ok) throw new Error('Failed to fetch metrics')

      const peersData = await peersRes.json()
      const metricsData = await metricsRes.json()

      setPeers(peersData.peers || [])
      setMetrics(metricsData)
      setLastUpdate(new Date())
      setLoading(false)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error')
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
    const interval = setInterval(fetchData, 5000)
    return () => clearInterval(interval)
  }, [])

  const isPeerActive = (lastSeen: string) => {
    const now = new Date()
    const seen = new Date(lastSeen)
    const diffMinutes = (now.getTime() - seen.getTime()) / 1000 / 60
    return diffMinutes < 5
  }

  const activePeers = peers.filter(p => isPeerActive(p.last_seen))

  return (
    <div className="flex min-h-screen bg-[#0a0a0a]">
      {/* Sidebar */}
      <div className="w-64 border-r border-neutral-800/10 flex flex-col bg-[#0a0a0a]">
        <div className="p-6 border-b border-neutral-800/10">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 rounded-lg bg-white flex items-center justify-center">
              <Network className="w-5 h-5 text-black" />
            </div>
            <span className="text-xl font-semibold text-white">ShadowNet</span>
          </div>
        </div>

        <nav className="flex-1 p-4 space-y-1">
          <a href="#" className="flex items-center gap-3 px-3 py-2 rounded-lg bg-neutral-900 text-white">
            <Home className="w-4 h-4" />
            <span className="text-sm">Dashboard</span>
          </a>
          <a href="#" className="flex items-center gap-3 px-3 py-2 rounded-lg text-neutral-400 hover:bg-neutral-900 hover:text-white transition-colors">
            <Network className="w-4 h-4" />
            <span className="text-sm">Peers</span>
          </a>
          <a href="#" className="flex items-center gap-3 px-3 py-2 rounded-lg text-neutral-400 hover:bg-neutral-900 hover:text-white transition-colors">
            <BarChart3 className="w-4 h-4" />
            <span className="text-sm">Analytics</span>
          </a>
          <a href="#" className="flex items-center gap-3 px-3 py-2 rounded-lg text-neutral-400 hover:bg-neutral-900 hover:text-white transition-colors">
            <Shield className="w-4 h-4" />
            <span className="text-sm">Security</span>
          </a>
        </nav>

        <div className="p-4 border-t border-neutral-800/10">
          <a href="#" className="flex items-center gap-3 px-3 py-2 rounded-lg text-neutral-400 hover:bg-neutral-900 hover:text-white transition-colors">
            <Settings className="w-4 h-4" />
            <span className="text-sm">Settings</span>
          </a>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 overflow-auto">
        {/* Top Bar */}
        <div className="bg-[#0a0a0a] sticky top-0 z-10">
          <div className="px-8 py-4 flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-semibold text-white">Dashboard</h1>
              <p className="text-sm text-neutral-500 mt-1">Monitor your P2P mesh VPN network</p>
            </div>
            <Button onClick={fetchData} variant="outline" size="sm">
              <RefreshCw className="w-4 h-4 mr-2" />
              Refresh
            </Button>
          </div>
        </div>

        <div className="p-8 space-y-6">
          {/* Error State */}
          {error && (
            <div className="bg-red-950/20 rounded-lg p-4">
              <div className="flex items-center gap-3 text-red-400">
                <WifiOff className="w-5 h-5" />
                <div>
                  <p className="font-medium text-sm">Connection Error</p>
                  <p className="text-xs text-red-300/70 mt-1">{error}</p>
                </div>
              </div>
            </div>
          )}

          {/* Metrics Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <Card>
              <CardHeader className="pb-3">
                <div className="flex items-center justify-between">
                  <CardDescription className="text-xs text-neutral-500">Active Peers</CardDescription>
                  <Activity className="w-4 h-4 text-green-500" />
                </div>
              </CardHeader>
              <CardContent>
                <div className="text-3xl font-bold text-white">
                  {loading ? '...' : activePeers.length}
                </div>
                <p className="text-xs text-neutral-600 mt-1">Currently connected</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="pb-3">
                <div className="flex items-center justify-between">
                  <CardDescription className="text-xs text-neutral-500">Total Peers</CardDescription>
                  <Users className="w-4 h-4 text-blue-500" />
                </div>
              </CardHeader>
              <CardContent>
                <div className="text-3xl font-bold text-white">
                  {loading ? '...' : peers.length}
                </div>
                <p className="text-xs text-neutral-600 mt-1">All registered</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="pb-3">
                <div className="flex items-center justify-between">
                  <CardDescription className="text-xs text-neutral-500">Control Plane</CardDescription>
                  <Server className="w-4 h-4 text-purple-500" />
                </div>
              </CardHeader>
              <CardContent>
                <div className="text-3xl font-bold text-white">
                  {loading ? '...' : error ? 'Offline' : 'Online'}
                </div>
                <p className="text-xs text-neutral-600 mt-1">{metrics?.uptime || 'N/A'}</p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="pb-3">
                <div className="flex items-center justify-between">
                  <CardDescription className="text-xs text-neutral-500">Network Health</CardDescription>
                  <Zap className="w-4 h-4 text-orange-500" />
                </div>
              </CardHeader>
              <CardContent>
                <div className="text-3xl font-bold text-white">
                  {loading ? '...' : peers.length > 0 ? Math.round((activePeers.length / peers.length) * 100) + '%' : '0%'}
                </div>
                <p className="text-xs text-neutral-600 mt-1">Uptime ratio</p>
              </CardContent>
            </Card>
          </div>

          {/* Network Topology */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Network Topology</CardTitle>
              <CardDescription className="text-sm">Real-time visualization of mesh network connections</CardDescription>
            </CardHeader>
            <CardContent className="p-0">
              {peers.length === 0 ? (
                <div className="flex flex-col items-center justify-center py-20 text-neutral-500">
                  <WifiOff className="w-12 h-12 mb-4 text-neutral-700" />
                  <p className="text-sm font-medium">No peers connected</p>
                  <p className="text-xs text-neutral-600 mt-1">Start a node to see it appear here</p>
                </div>
              ) : (
                <NetworkTopology peers={peers} />
              )}
            </CardContent>
          </Card>

          {/* Peers Table */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Connected Peers</CardTitle>
              <CardDescription className="text-sm">Detailed status of all peers in the network</CardDescription>
            </CardHeader>
            <CardContent>
              {loading && peers.length === 0 ? (
                <div className="text-center py-8 text-neutral-500 text-sm">
                  Loading peers...
                </div>
              ) : peers.length === 0 ? (
                <div className="text-center py-8">
                  <p className="text-neutral-500 text-sm">No peers connected</p>
                </div>
              ) : (
                <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="">
                        <th className="text-left py-3 px-4 font-medium text-neutral-500 text-xs">Status</th>
                        <th className="text-left py-3 px-4 font-medium text-neutral-500 text-xs">Peer ID</th>
                        <th className="text-left py-3 px-4 font-medium text-neutral-500 text-xs">Endpoint</th>
                        <th className="text-left py-3 px-4 font-medium text-neutral-500 text-xs">Public Key</th>
                        <th className="text-left py-3 px-4 font-medium text-neutral-500 text-xs">Last Seen</th>
                      </tr>
                    </thead>
                    <tbody>
                      {peers.map((peer) => {
                        const active = isPeerActive(peer.last_seen)
                        return (
                          <tr key={peer.id} className="hover:bg-neutral-900/30 transition-colors">
                            <td className="py-3 px-4">
                              <div className="flex items-center gap-2">
                                <div className={`w-2 h-2 rounded-full ${active ? 'bg-green-500' : 'bg-neutral-600'}`} />
                                <Badge variant={active ? "success" : "secondary"} className="text-xs">
                                  {active ? 'Active' : 'Inactive'}
                                </Badge>
                              </div>
                            </td>
                            <td className="py-3 px-4">
                              <code className="text-xs text-white font-mono bg-neutral-900 px-2 py-1 rounded">
                                {peer.id}
                              </code>
                            </td>
                            <td className="py-3 px-4">
                              <span className="text-xs text-neutral-300">
                                {peer.endpoint_ip}:{peer.endpoint_port}
                              </span>
                            </td>
                            <td className="py-3 px-4">
                              <code className="text-xs text-neutral-500 font-mono">
                                {peer.wg_public_key.substring(0, 20)}...
                              </code>
                            </td>
                            <td className="py-3 px-4 text-xs text-neutral-500">
                              {timeAgo(peer.last_seen)}
                            </td>
                          </tr>
                        )
                      })}
                    </tbody>
                  </table>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
