dev:
	air -c air-dev.toml

stag:
	air -c air-stag.toml

hand-deploy-to-stag-step-1:
	git push -u origin develop; \
	git checkout stag-release; \
	git pull origin develop;

hand-deploy-to-stag-step-2:
	git push -u origin stag-release; \
	git checkout develop;