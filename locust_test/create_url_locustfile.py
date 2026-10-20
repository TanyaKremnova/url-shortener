# @task
# def create_url(self):
#     with self.client.post(
#         "/urls/",
#         json={"original_url": "https://www.github.com"},
#         headers={"Authorization": f"Bearer {self.token}"},
#         name="/urls/ [create]",
#         catch_response=True,
#     ) as response:
#         if response.status_code != 201:
#             response.failure(
#                 f"Expected 201, got {response.status_code}"
#             )


# locust -f ./locust_test/create_url_locustfile.py \
#   --host http://localhost:8081 \
#   --users 50 \
#   --spawn-rate 5 \
#   --run-time 3m \
#   --headless \
#   --csv=results/app1-create
