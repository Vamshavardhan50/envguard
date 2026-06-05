# testdata/python-project/app.py
# Sample Python usage for envguard tests.

import os

debug = os.getenv("DEBUG")
key = os.environ.get("SECRET_KEY")
missing = os.getenv("MISSING_VAR")

_ = debug
_ = key
_ = missing
