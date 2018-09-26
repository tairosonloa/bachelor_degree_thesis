# -*- coding: utf-8 -*-

import json


CONFIG_FILE = "../../config.json" # File with API IP, API port, Bearer token and more


def load_config():
    """Loads config from JSON file and return its content as a dictionary"""
    with open(CONFIG_FILE, "r") as f:
        config = json.load(f)
    return config