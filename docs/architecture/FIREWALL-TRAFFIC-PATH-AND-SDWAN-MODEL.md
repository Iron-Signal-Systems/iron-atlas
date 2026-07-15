# Firewall Traffic Path and SD-WAN Model

## Path Evaluation Input

- Firewall and routing domain or VDOM
- Source interface, zone, address, identity, and device context
- Destination address or object
- Protocol and ports
- Time and schedule context
- Application or internet-service context when available

## Evaluation

1. Resolve ingress interface and routing domain.
2. Expand source, destination, service, schedule, and nested groups.
3. Preserve vendor rule order and disabled state.
4. Identify the first applicable policy or vendor-equivalent rule.
5. Evaluate policy routing or rule-selected gateway behavior.
6. Evaluate SD-WAN rule order, match criteria, strategy, health requirements, preference, and fallback.
7. Evaluate connected and longest-prefix route candidates, distance, priority, and next-hop reachability.
8. Resolve physical, logical, tunnel, or SD-WAN egress.
9. Resolve source and destination translation.
10. Resolve inspection, logging, shaping, and authentication controls.
11. Evaluate known return-path asymmetry risks.
12. Attach evidence state and source location to each conclusion.

## Vendor Differences

FortiGate may reference SD-WAN zones from routes and policies and use separate SD-WAN rules for member selection. OPNsense and pfSense commonly select gateways or gateway groups from interface firewall rules. The normalized model expresses common intent without erasing those vendor-specific semantics.

## Result

The user receives a human-readable explanation and a structured path record. Runtime health-dependent conclusions remain unknown unless current operational evidence supports them.
