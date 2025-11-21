import { tsParticles } from "@tsparticles/engine";
import { loadSlim } from "@tsparticles/slim";

export async function initParticles() {
    await loadSlim(tsParticles);

    // Responsive particle count based on screen size
    const isMobile = window.innerWidth < 768;
    const isTablet = window.innerWidth >= 768 && window.innerWidth < 1024;

    let particleCount = 80; // Desktop default
    if (isMobile) {
        particleCount = 25; // Mobile: much fewer
    } else if (isTablet) {
        particleCount = 45; // Tablet: moderate
    }

    await tsParticles.load({
        id: "tsparticles",
        options: {
            particles: {
                destroy: {
                    mode: "split",
                    split: {
                        count: 1,
                        factor: {
                            value: {
                                min: 2,
                                max: 4
                            }
                        },
                        rate: {
                            value: 100
                        },
                        particles: {
                            life: {
                                count: 1,
                                duration: {
                                    value: {
                                        min: 2,
                                        max: 3
                                    }
                                }
                            },
                            move: {
                                speed: {
                                    min: 1,
                                    max: 5
                                }
                            }
                        }
                    }
                },
                number: {
                    value: particleCount
                },
                color: {
                    value: [
                        "#10b981", // emerald - doppler
                        "#14b8a6", // teal - luna
                        "#3b82f6", // blue - augur
                        "#06b6d4", // cyan - lunar prance
                        "#ec4899", // pink - kami
                        "#8b5cf6", // violet - dreamflow
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
                    value: 0.95,
                    animation: {
                        enable: true,
                        speed: 0.3,
                        minimumValue: 0.6
                    }
                },
                size: {
                    value: {
                        min: 18,
                        max: 28
                    }
                },
                collisions: {
                    enable: !isMobile, // Disable collisions on mobile for better performance
                    mode: "bounce"
                },
                move: {
                    enable: true,
                    speed: 1.5,
                    direction: "none",
                    random: true,
                    straight: false,
                    outModes: "bounce",
                    attract: {
                        enable: false
                    }
                }
            },
            interactivity: {
                events: {
                    onClick: {
                        enable: true,
                        mode: "push"
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

