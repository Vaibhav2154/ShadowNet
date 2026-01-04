'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Home, ChevronRight, BookOpen } from 'lucide-react'
import Link from 'next/link'
import { DOCS_SECTIONS, SectionId } from '@/lib/docs-data'

export default function DocsPage() {
  const [activeSection, setActiveSection] = useState<SectionId>('intro')

  const ActiveComponent = DOCS_SECTIONS.find(s => s.id === activeSection)?.component || (() => null)

  return (
    <div className="flex min-h-screen bg-black">
      {/* Sidebar */}
      <div className="w-64 border-r border-neutral-800/50 p-6 space-y-6 hidden md:block sticky top-0 h-screen overflow-y-auto">
        <div className="space-y-3">
          <Link href="/">
            <Button variant="ghost" size="sm" className="w-full justify-start pl-0 hover:bg-transparent hover:text-white text-neutral-400">
              <Home className="w-4 h-4 mr-2" />
              Back to Dashboard
            </Button>
          </Link>

          <div className="flex items-center gap-2 py-2">
            <BookOpen className="w-5 h-5 text-white" />
            <h2 className="text-lg font-semibold text-white">Documentation</h2>
          </div>
        </div>

        <nav className="space-y-1">
          {DOCS_SECTIONS.map((section) => {
            const Icon = section.icon
            return (
              <button
                key={section.id}
                onClick={() => setActiveSection(section.id)}
                className={`w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors ${activeSection === section.id
                    ? 'bg-neutral-900 text-white font-medium shadow-sm border border-neutral-800'
                    : 'text-neutral-400 hover:text-white hover:bg-neutral-900/50'
                  }`}
              >
                <Icon className="w-4 h-4" />
                <span className="flex-1 text-left">{section.title}</span>
                {activeSection === section.id && <ChevronRight className="w-4 h-4 text-neutral-500" />}
              </button>
            )
          })}
        </nav>
      </div>

      {/* Main Content */}
      <div className="flex-1 overflow-auto">
        <div className="max-w-5xl mx-auto p-8 space-y-8 pb-32">
          <ActiveComponent />
        </div>
      </div>
    </div>
  )
}
