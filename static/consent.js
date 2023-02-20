window.addEventListener('load', function () {
	const cc = initCookieConsent();
	cc.run({
		current_lang: 'es',
		page_scripts: true,
		mode: 'opt-in',
		cookie_name: 'consented',
		cookie_expiration: 360,

		onFirstAction: function (user_preferences, cookie) {
		},

		onAccept: function (cookie) {
			gtag('consent', 'update', {
				'analytics_storage': 'granted'
			});
			window.clarity('consent');
		},

		onChange: function (cookie, changed_categories) {
		},

		languages: {
			'es': {
				consent_modal: {
					title: 'Utilizamos cookies',
					description: 'Vigo360 usa Google Analytics con fines estadísticos. Si pulsas aceptar, aceptas el uso de cookies de terceros y que se compartan dichos datos con Google. Rechazarlo ocultará este mensaje y no se compartirán datos con terceros..',
					primary_btn: {
						text: 'Aceptar',
						role: 'accept_all'
					},
					secondary_btn: {
						text: 'Rechazar',
						role: 'accept_necessary'
					}
				},
				settings_modal: {
					title: 'Cookie preferences',
					save_settings_btn: 'Save settings',
					accept_all_btn: 'Accept all',
					reject_all_btn: 'Reject all',
					close_btn_label: 'Close',
					cookie_table_headers: [
						{col1: 'Name'},
						{col2: 'Domain'},
						{col3: 'Expiration'},
						{col4: 'Description'}
					],
					blocks: [
						{
							title: 'Cookie usage 📢',
							description: 'I use cookies to ensure the basic functionalities of the website and to enhance your online experience. You can choose for each category to opt-in/out whenever you want. For more details relative to cookies and other sensitive data, please read the full <a href="#" class="cc-link">privacy policy</a>.'
						}, {
							title: 'Strictly necessary cookies',
							description: 'These cookies are essential for the proper functioning of my website. Without these cookies, the website would not work properly',
							toggle: {
								value: 'necessary',
								enabled: true,
								readonly: true          // cookie categories with readonly=true are all treated as "necessary cookies"
							}
						}, {
							title: 'Performance and Analytics cookies',
							description: 'These cookies allow the website to remember the choices you have made in the past',
							toggle: {
								value: 'analytics',
								enabled: false,
								readonly: false
							},
							cookie_table: [
								{
									col1: '^_ga',
									col2: 'google.com',
									col3: '2 years',
									col4: 'description ...',
									is_regex: true
								},
								{
									col1: '_gid',
									col2: 'google.com',
									col3: '1 day',
									col4: 'description ...',
								}
							]
						}
					]
				}

			}
		}
	});
});