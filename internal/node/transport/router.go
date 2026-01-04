package transport

import (
	"log"
)

// Router handles packet routing between TUN and WireGuard
// In userspace WireGuard, routing is handled by the WireGuard device itself
// This file is a placeholder for any custom routing logic

// Router manages packet routing
type Router struct {
	// Future: custom routing tables, policy routing, etc.
}

// NewRouter creates a new router
func NewRouter() *Router {
	return &Router{}
}

// Start starts the router
func (r *Router) Start() error {
	log.Println("Router started (routing handled by WireGuard device)")
	return nil
}

// Stop stops the router
func (r *Router) Stop() error {
	log.Println("Router stopped")
	return nil
}

// AddRoute adds a route (placeholder)
func (r *Router) AddRoute(destination, gateway string) error {
	// Custom routing logic could be added here
	return nil
}

// RemoveRoute removes a route (placeholder)
func (r *Router) RemoveRoute(destination string) error {
	// Custom routing logic could be added here
	return nil
}
