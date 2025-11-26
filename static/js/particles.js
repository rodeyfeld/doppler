import { tsParticles } from "@tsparticles/engine";
import { loadSlim } from "@tsparticles/slim";
import { loadImageShape } from "@tsparticles/shape-image";

export async function initParticles() {
    await loadSlim(tsParticles);
    await loadImageShape(tsParticles);

    // Responsive particle count based on screen size
    const isMobile = window.innerWidth < 768;
    const isTablet = window.innerWidth >= 768 && window.innerWidth < 1024;

    let particleCount = 60;
    if (isMobile) {
        particleCount = 15;
    } else if (isTablet) {
        particleCount = 30;
    }

    await tsParticles.load({
        id: "tsparticles",
        options: {
            particles: {
                number: {
                    value: particleCount,
                    density: {
                        enable: true,
                        area: 800
                    }
                },
                color: {
                    value: [
                        "#1eb854", 
                        "#1db88e", 
                        "#1db8ab", 
                        "#38bdf8", 
                        "#36d399", 
                        "#fbbd23", 
                        "#f87272", 
                        "#f5f5f4", 
                        "#3f4c46", 
                    ]
                },
                rotate: {
                    direction: "random",
                    value: { min: 0, max: 360 },
                    animation: {
                        speed: 25,
                        enable: true,
                    },
                },
                shape: {
                    type: "images",
                    options: {
                        images: [
                            {
                                src: "static/pinwheel.svg",
                                width: 100,
                                height: 100,
                                replaceColor: true
                            }
                        ]
                    }
                },
                opacity: {
                    value: 0.65,
                    random: {
                        enable: true,
                        minimumValue: 0.3
                    },
                    animation: {
                        enable: true,
                        speed: 0.12,
                        startValue: "random",
                        sync: false,
                        minimumValue: 0.3
                    }
                },
                size: {
                    value: {
                        min: 12,
                        max: 26
                    }
                },
                move: {
                    enable: true,
                    speed: {
                        min: 0.05,
                        max: 0.25
                    },
                    direction: "top",
                    random: true,
                    drift: 0,
                    straight: false,
                    outModes: {
                        default: "out"
                    },
                    attract: {
                        enable: false
                    }
                }
            },
            interactivity: {
                events: {
                    onClick: {
                        enable: false
                    },
                    onHover: {
                        enable: false
                    }
                },
                modes: {
                    push: {
                        quantity: 3
                    }
                }
            },
        }
    });
}

