# apixen
This is an example of a REST API that uses Go, Gorilla Mux and a PostgreSQL database. It uses Docker Compose to run multiple-container Docker applications.

Prerequisite:-
- Docker
- Docker Compose

Installation:-
1. Install Docker. You can find instructions for the installation here; https://docs.docker.com/install/ .
2. Install Docker Compose. You can find instructions for the installation here; https://docs.docker.com/compose/install/.
4. On terminal, git clone https://github.com/maszuari/apixen.git
3. Run this command: docker-compose up

Use these commands:-
1. To get members of an organization: curl -X GET http://localhost:3000/orgs/org-name/members/
- For example: curl -X GET http://localhost:3000/orgs/acme/members/
2. To save comment that associate with an organizations: curl -X POST http://localhost:3000/orgs/org-name/comments/ -H 'Content-Type: application/json' -d '{"comment": "Write comment here"}'
- For example: curl -X POST http://localhost:3000/orgs/acme/comments/ -H 'Content-Type: application/json' -d '{"comment": "Comment 123"}'
3. To get all comments of an organization: curl -X GET http://localhost:3000/orgs/org-name/comments/
- For example: curl -X GET http://localhost:3000/orgs/acme/comments/
4. To delete all comments that associate with an organization: curl -X DELETE http://localhost:3000/orgs/org-name/comments
- For example curl -X DELETE http://localhost:3000/orgs/acme/comments
- In the URL, there is no '/' after 'comments'
