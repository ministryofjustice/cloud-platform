PUBLISHING_IMAGE := ministryofjustice/cloud-platform-tech-docs-publisher:1.6

# Use this to run a local instance of the documentation site, while editing
.PHONY: preview
preview:
	docker run --rm \
		-v $$(pwd)/config.rb:/app/config.rb \
		-v $$(pwd)/config:/app/config \
		-v $$(pwd)/source:/app/source \
		-p 4567:4567 \
		$(PUBLISHING_IMAGE) \
		bundle exec middleman serve

# Build the HTML files and run the broken link checker
.PHONY: check
check:
	make build
	make htmlproofer

# Build the docs/ files locally
.PHONY: build
build:
	docker run --rm \
		-v $$(pwd)/config.rb:/app/config.rb \
		-v $$(pwd)/config:/app/config \
		-v $$(pwd)/source:/app/source \
		-v $$(pwd)/docs:/app/docs \
		-p 4567:4567 \
		$(PUBLISHING_IMAGE) \
		bundle exec middleman build --build-dir docs 2>&1 | grep -v 'warning: URI.*escape is obsolete'

# Check for broken links (assumes `make build` has been run)
.PHONY: htmlproofer
htmlproofer:
	docker run --rm \
		-v $$(pwd)/config.rb:/app/config.rb \
		-v $$(pwd)/config:/app/config \
		-v $$(pwd)/source:/app/source \
		-v $$(pwd)/docs:/app/docs \
		-p 4567:4567 \
		$(PUBLISHING_IMAGE) \
		bundle exec htmlproofer ./docs --http-status-ignore 429 --allow-hash-href --http-status-ignore 0,429 --url-swap "https?\:\/\/user-guide\.cloud-platform\.service\.justice\.gov\.uk:" ./docs


