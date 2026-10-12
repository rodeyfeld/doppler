import { tsParticles } from "@tsparticles/engine";
import { loadSlim } from "@tsparticles/slim";
import { loadImageShape } from "@tsparticles/shape-image";

// A sparse field of small pinwheels in the site's greens and blues, drifting
// slowly across the landing page behind the text.
export async function initParticles() {
    await loadSlim(tsParticles);
    await loadImageShape(tsParticles);

    const width = window.innerWidth;
    const particleCount = width < 768 ? 24 : width < 1024 ? 40 : 70;

    await tsParticles.load({
        id: "tsparticles",
        options: {
            fpsLimit: 30,
            detectRetina: true,
            particles: {
                number: {
                    value: particleCount,
                    density: { enable: true, area: 900 },
                },
                color: {
                    value: [
                        "#14b8a6", "#06b6d4", "#10b981", "#0ea5e9", "#22c55e",
                        "#22d3ee", "#3b82f6", "#34d399", "#0d9488", "#38bdf8",
                    ],
                },
                // Each pinwheel gets its own spin: some barely turn, some whirl
                rotate: {
                    direction: "random",
                    value: { min: 0, max: 360 },
                    animation: { enable: true, speed: { min: 1, max: 14 }, sync: false },
                },
                shape: {
                    type: "images",
                    options: {
                        images: [
                            {
                                src: "static/pinwheel.svg",
                                width: 100,
                                height: 100,
                                replaceColor: true,
                            },
                        ],
                    },
                },
                opacity: {
                    value: { min: 0.18, max: 0.45 },
                    animation: {
                        enable: true,
                        speed: 0.15,
                        startValue: "random",
                        sync: false,
                    },
                },
                size: { value: { min: 7, max: 16 } },
                // ...and its own pace across the page, from a crawl to a breeze
                move: {
                    enable: true,
                    speed: { min: 0.05, max: 0.7 },
                    direction: "right",
                    straight: true,
                    outModes: { default: "out" },
                },
            },
            interactivity: {
                events: {
                    onClick: { enable: false },
                    onHover: { enable: false },
                },
            },
        },
    });
}
