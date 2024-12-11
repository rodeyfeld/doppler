
//tsParticles library - https://github.com/matteobruni/tsparticles

//tsParticles library - https://github.com/matteobruni/tsparticles

(async () => {
	await loadFull(tsParticles);

	await tsParticles.load({
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
					value: 80
				},
				color: {
					value: [
						"#1eb854",
						"#1db88e",
						"#1db8ab",
						"#19362d",
					]
				},
				rotate: {
					direction: "random",
					value: { min: 0, max: 360 },
					animation: { speed: 20, enable: true, },
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
					value: 1
				},
				size: {
					value: {
						min: 10,
						max: 15
					}
				},
				collisions: {
					enable: true,
					mode: "bounce"
				},
				move: {
					enable: true,
					speed: 1,
					outModes: "bounce"
				}
			},
			interactivity: {
				events: {
					onClick: {
						enable: true,
						mode: "pop"
					}
				}
			},
		}
	});
})();
