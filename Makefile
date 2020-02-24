.PHONY: deploy destroy

deploy:
	cd cdk && npm install && npm run build && npm run cdk -- deploy

destroy:
	cd cdk && npm install && npm run cdk -- destroy