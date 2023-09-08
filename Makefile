ts = $(shell date +%s)

docker:
	docker build . \
	 	-t shadowsocks/$(app):$(ts)\
	 	-t shadowsocks/$(app):latest\
		--build-arg app=$(app)
