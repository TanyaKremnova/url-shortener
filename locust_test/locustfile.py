from locust import HttpUser, task, between
import random
import string

def random_code():
    return ''.join(random.choices(string.ascii_letters, k=6))

class URLShortenerUser(HttpUser):
    wait_time = between(0.1, 0.5)
    token = None
    short_codes = []

    def on_start(self):
        # Register and login once per simulated user
        email = f"user_{random_code()}@test.com"
        self.client.post("/auth/register", json={
            "email": email,
            "password": "password123"
        })
        response = self.client.post("/auth/login", json={
            "email": email,
            "password": "password123"
        })
        data = response.json()
        self.token = data.get("data", {}).get("token")

    @task(1)
    def create_url(self):
        if not self.token:
            return
        try:
            response = self.client.post("/urls/",
                json={"original_url": "https://www.github.com"},
                headers={"Authorization": f"Bearer {self.token}"}
            )
            if response.status_code == 201:
                data = response.json()
                code = data.get("data", {}).get("short_code")
                if code:
                    self.short_codes.append(code)
        except Exception:
            pass

    @task(10)
    def redirect(self):
        # Redirects are 10x more frequent than creates — realistic ratio
        if not self.short_codes:
            return
        code = random.choice(self.short_codes)
        self.client.get(f"/{code}", allow_redirects=False)

    @task(2)
    def get_stats(self):
        if not self.token:
            return
        self.client.get("/admin/urls/stats",
            headers={"Authorization": f"Bearer {self.token}"}
        )