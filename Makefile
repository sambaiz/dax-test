.PHONY: deploy

deploy:
	cd cdk && npm run build && npm run cdk -- deploy