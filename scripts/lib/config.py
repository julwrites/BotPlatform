import os
import sys
from scripts.lib.yaml import SimpleYaml

# Defaults mirroring the old hardcoded values
DEFAULT_CONFIG = {
    "project": {
        "name": "AgentHarness",
        "root": "."
    },
    "tasks": {
        "dir": "docs/tasks",
        "categories": [
            "foundation",
            "infrastructure",
            "domain",
            "presentation",
            "migration",
            "features",
            "testing",
            "review",
            "security",
            "research"
        ],
        "statuses": [
            "pending",
            "in_progress",
            "wip_blocked",
            "review_requested",
            "verified",
            "completed",
            "blocked",
            "cancelled",
            "deferred"
        ],
        "types": [
            "epic",
            "story",
            "task",
            "bug"
        ]
    },
    "agents": {
        "bus_dir": ".agents",
        "audit_log": "logs/audit.jsonl"
    }
}

class Config:
    _instance = None

    @staticmethod
    def load(repo_root=None):
        if Config._instance:
            return Config._instance

        if not repo_root:
            # Try to deduce repo root
            script_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__))) # scripts/
            repo_root = os.getenv("TASKS_REPO_ROOT", os.path.dirname(script_dir))

        config_path = os.path.join(repo_root, "harness.config.yaml")

        user_config = {}
        if os.path.exists(config_path):
            try:
                user_config = SimpleYaml.load(config_path)
            except Exception as e:
                print(f"Warning: Failed to load harness.config.yaml: {e}")

        # Deep merge with defaults
        Config._instance = Config._merge(DEFAULT_CONFIG, user_config)
        Config._instance["_root"] = repo_root
        return Config._instance

    @staticmethod
    def _merge(default, user):
        result = default.copy()
        for key, val in user.items():
            if key in result and isinstance(result[key], dict) and isinstance(val, dict):
                result[key] = Config._merge(result[key], val)
            else:
                result[key] = val
        return result

def get_config(repo_root=None):
    return Config.load(repo_root)
