'use client'

import { useEffect, useRef } from 'react'

interface Node {
  id: string
  x: number
  y: number
  status: 'active' | 'inactive'
}

interface NetworkTopologyProps {
  peers: Array<{
    id: string
    last_seen: string
  }>
}

export function NetworkTopology({ peers }: NetworkTopologyProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null)

  const isPeerActive = (lastSeen: string) => {
    const now = new Date()
    const seen = new Date(lastSeen)
    const diffMinutes = (now.getTime() - seen.getTime()) / 1000 / 60
    return diffMinutes < 5
  }

  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas) return

    const ctx = canvas.getContext('2d')
    if (!ctx) return

    // Set canvas size properly
    const dpr = window.devicePixelRatio || 1
    const rect = canvas.getBoundingClientRect()

    const width = rect.width || 800
    const height = rect.height || 500

    canvas.width = width * dpr
    canvas.height = height * dpr
    ctx.scale(dpr, dpr)

    // Clear canvas with black background
    ctx.fillStyle = '#000000'
    ctx.fillRect(0, 0, width, height)

    // Control plane node (center)
    const centerX = width / 2
    const centerY = height / 2
    const controlPlane = { x: centerX, y: centerY, id: 'Control Plane' }

    // Position peer nodes in a large circle
    const minDimension = Math.min(width, height)
    const radius = minDimension * 0.38

    const nodes: Node[] = peers.map((peer, i) => {
      const angle = (i / peers.length) * 2 * Math.PI - Math.PI / 2
      const x = centerX + radius * Math.cos(angle)
      const y = centerY + radius * Math.sin(angle)

      return {
        id: peer.id,
        x,
        y,
        status: isPeerActive(peer.last_seen) ? 'active' : 'inactive'
      }
    })

    // Draw P2P connections (mesh) - only between active peers
    const activeNodes = nodes.filter(n => n.status === 'active')
    activeNodes.forEach((node, i) => {
      activeNodes.slice(i + 1).forEach(otherNode => {
        ctx.beginPath()
        ctx.moveTo(node.x, node.y)
        ctx.lineTo(otherNode.x, otherNode.y)
        ctx.strokeStyle = 'rgba(34, 197, 94, 0.1)'
        ctx.lineWidth = 1
        ctx.setLineDash([3, 3])
        ctx.stroke()
        ctx.setLineDash([])
      })
    })

    // Draw connections from control plane to peers
    nodes.forEach(node => {
      ctx.beginPath()
      ctx.moveTo(controlPlane.x, controlPlane.y)
      ctx.lineTo(node.x, node.y)
      ctx.strokeStyle = node.status === 'active' ? 'rgba(139, 92, 246, 0.2)' : 'rgba(64, 64, 64, 0.3)'
      ctx.lineWidth = node.status === 'active' ? 1.5 : 1
      ctx.stroke()
    })

    // Draw control plane node
    ctx.beginPath()
    ctx.arc(controlPlane.x, controlPlane.y, 20, 0, 2 * Math.PI)
    ctx.fillStyle = '#8b5cf6'
    ctx.fill()
    ctx.strokeStyle = '#a78bfa'
    ctx.lineWidth = 2
    ctx.stroke()

    // Control plane glow
    ctx.beginPath()
    ctx.arc(controlPlane.x, controlPlane.y, 28, 0, 2 * Math.PI)
    ctx.strokeStyle = 'rgba(139, 92, 246, 0.2)'
    ctx.lineWidth = 1
    ctx.stroke()

    // Draw peer nodes
    nodes.forEach(node => {
      // Node circle
      ctx.beginPath()
      ctx.arc(node.x, node.y, 14, 0, 2 * Math.PI)
      ctx.fillStyle = node.status === 'active' ? '#22c55e' : '#404040'
      ctx.fill()
      ctx.strokeStyle = node.status === 'active' ? '#4ade80' : '#525252'
      ctx.lineWidth = 2
      ctx.stroke()

      // Pulse effect for active nodes
      if (node.status === 'active') {
        ctx.beginPath()
        ctx.arc(node.x, node.y, 20, 0, 2 * Math.PI)
        ctx.strokeStyle = 'rgba(34, 197, 94, 0.3)'
        ctx.lineWidth = 1
        ctx.stroke()
      }
    })

    // Draw labels with better positioning
    ctx.textAlign = 'center'
    ctx.textBaseline = 'middle'

    // Control plane label
    ctx.fillStyle = '#ffffff'
    ctx.font = 'bold 12px sans-serif'
    ctx.fillText('Control Plane', controlPlane.x, controlPlane.y + 42)

    // Peer labels - FULL NAMES
    nodes.forEach(node => {
      ctx.fillStyle = node.status === 'active' ? '#ffffff' : '#737373'
      ctx.font = node.status === 'active' ? '11px monospace' : '10px monospace'

      // Draw full peer ID
      ctx.fillText(node.id, node.x, node.y + 32)
    })

  }, [peers])

  return (
    <div className="relative w-full h-[500px] bg-[#0a0a0a] rounded-lg">
      <canvas
        ref={canvasRef}
        className="w-full h-full"
      />

      {/* Legend */}
      <div className="absolute bottom-4 right-4 bg-[#0f0f0f] border border-neutral-800/10 rounded-lg px-4 py-2">
        <div className="flex items-center gap-4 text-xs">
          <div className="flex items-center gap-2">
            <div className="w-2.5 h-2.5 rounded-full bg-purple-500" />
            <span className="text-neutral-400">Control Plane</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-2.5 h-2.5 rounded-full bg-green-500" />
            <span className="text-neutral-400">Active</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-2.5 h-2.5 rounded-full bg-neutral-600" />
            <span className="text-neutral-400">Inactive</span>
          </div>
        </div>
      </div>
    </div>
  )
}
