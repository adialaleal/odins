package caddy

// buildBaseConfig returns the initial Caddy JSON config with TLS internal.
func buildBaseConfig(tld string) map[string]interface{} {
	return map[string]interface{}{
		"apps": map[string]interface{}{
			"tls": map[string]interface{}{
				"automation": map[string]interface{}{
					"policies": []interface{}{
						map[string]interface{}{
							"issuers": []interface{}{
								map[string]interface{}{"module": "internal"},
							},
						},
					},
				},
			},
			"http": map[string]interface{}{
				"servers": map[string]interface{}{
					"srv0": map[string]interface{}{
						"listen": []string{":443", ":80"},
						"routes": []interface{}{},
					},
				},
			},
		},
	}
}

// buildDomainRoute returns a Caddy route that serves a static file_server for a domain landing page.
func buildDomainRoute(hostname, pageDir string) map[string]interface{} {
	return map[string]interface{}{
		"@id": "odins-domain-" + hostname,
		"match": []interface{}{
			map[string]interface{}{
				"host": []string{hostname},
			},
		},
		"handle": []interface{}{
			map[string]interface{}{
				"handler": "subroute",
				"routes": []interface{}{
					map[string]interface{}{
						"handle": []interface{}{
							map[string]interface{}{
								"handler": "file_server",
								"root":    pageDir,
							},
						},
					},
				},
			},
		},
		"terminal": true,
	}
}

// buildRoute returns a Caddy route JSON object for a single subdomain.
func buildRoute(subdomain, upstream, id string) map[string]interface{} {
	return map[string]interface{}{
		"@id": id,
		"match": []interface{}{
			map[string]interface{}{
				"host": []string{subdomain},
			},
		},
		"handle": []interface{}{
			map[string]interface{}{
				"handler": "subroute",
				"routes": []interface{}{
					map[string]interface{}{
						"handle": []interface{}{
							map[string]interface{}{
								"handler": "reverse_proxy",
								"upstreams": []interface{}{
									map[string]interface{}{
										"dial": upstream,
									},
								},
								"headers": map[string]interface{}{
									"request": map[string]interface{}{
										"set": map[string]interface{}{
											"X-Real-IP":       []string{"{http.request.remote.host}"},
											"X-Forwarded-For": []string{"{http.request.remote.host}"},
											"X-Forwarded-Proto": []string{"{http.request.scheme}"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"terminal": true,
	}
}
