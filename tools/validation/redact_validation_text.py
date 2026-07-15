#!/usr/bin/env python3
from pathlib import Path
import re
import sys

if len(sys.argv) != 3:
    raise SystemExit("usage: redact_validation_text.py INPUT OUTPUT")
text = Path(sys.argv[1]).read_text(errors="replace")
text = re.sub(r"(?i)(postgres(?:ql)?://)[^\s:@/]+:[^@\s]+@", r"\1REDACTED@", text)
text = re.sub(r"(?im)^(authorization:\s*)(?:bearer|basic)\s+\S+", r"\1REDACTED", text)
text = re.sub(r"(?im)^((?:password|token|secret|api[_-]?key)\s*[=:]\s*)\S+", r"\1REDACTED", text)
text = re.sub(r"-----BEGIN [^-]*PRIVATE KEY-----.*?-----END [^-]*PRIVATE KEY-----", "[REDACTED PRIVATE KEY]", text, flags=re.S)
Path(sys.argv[2]).write_text(text)
