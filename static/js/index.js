/**
 * Doppler - Main JavaScript Bundle
 * 
 * This bundles all JavaScript dependencies into a single file:
 * - HTMX for dynamic interactions
 * - Quill for rich text editing
 * - tsParticles for home page animations
 * 
 * Bun bundles everything from node_modules into static/dist/index.js
 * 
 * Build commands:
 *   - Development: bun run dev:js (watch mode)
 *   - Production: bun run build:js (minified)
 */

import htmx from 'htmx.org';
import { initParticles } from './particles.js';

// Make htmx globally available
window.htmx = htmx;

// Initialize components when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    // Initialize Quill editor for post creation
    const modal = document.getElementById('create_post_modal');
    if (modal) {
        import('./editor.js').then(({ initQuillEditor }) => {
            initQuillEditor();
        });
    }

    // Initialize particles on home page if container exists
    if (document.getElementById('tsparticles')) {
        initParticles();
    }
});
