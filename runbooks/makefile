IMAGE := ministryofjustice/cloud-platform-runbooks

server:
	mkdir compiled || true
	docker run \
		-p 4567:4567 \
		-v $$(pwd)/source:/app/source \
		-v $$(pwd)/compiled:/app/compiled \
		-v $$(pwd)/config:/app/config \
		-it \
		$(IMAGE) bundle exec middleman server

# The CircleCI build pipeline does this, so it should never
# be necessary to run this task. I'm leaving it here for
# reference, and in case we ever need to push changes
# manually.
build:
	mkdir compiled || true
	docker run \
		-v $$(pwd)/source:/app/source \
		-v $$(pwd)/compiled:/app/compiled \
		-v $$(pwd)/config:/app/config \
		-it \
		$(IMAGE) bundle exec middleman build --build-dir compiled
	touch compiled/.nojekyll
