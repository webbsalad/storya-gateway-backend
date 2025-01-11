PROTOGEN_DIR := vendor.protogen

TMP_DIR := ./tmp-storya-backend

API_DIRS := ./api/passport ./api/otp ./api/content ./api/recs

INTERNAL_PB_DIR := ./internal/pb


proto-deps:
	rm -rf $(PROTOGEN_DIR)
	mkdir -p $(PROTOGEN_DIR)

	git clone https://github.com/googleapis/googleapis.git $(PROTOGEN_DIR)/googleapis
	mv $(PROTOGEN_DIR)/googleapis/google/ $(PROTOGEN_DIR)
	rm -rf $(PROTOGEN_DIR)/googleapis/

	git clone --depth 1 https://github.com/bufbuild/protoc-gen-validate.git $(PROTOGEN_DIR)/protoc-gen-validate
	mv $(PROTOGEN_DIR)/protoc-gen-validate/validate/ $(PROTOGEN_DIR)
	rm -rf $(PROTOGEN_DIR)/protoc-gen-validate/


	for api_dir in $(API_DIRS); do \
		rm -rf $$api_dir; \
		mkdir -p $$api_dir; \
	done

	# passport
	git clone --depth 1 git@github.com:webbsalad/storya-passport-backend.git $(TMP_DIR)
	mv $(TMP_DIR)/api/passport api/
	rm -rf $(TMP_DIR)

	# otp
	git clone --depth 1 git@github.com:webbsalad/storya-otp-backend.git $(TMP_DIR)
	mv $(TMP_DIR)/api/otp api/
	rm -rf $(TMP_DIR)

	# content
	git clone --depth 1 git@github.com:webbsalad/storya-content-backend.git $(TMP_DIR)
	mv $(TMP_DIR)/api/content api/
	rm -rf $(TMP_DIR)

	# recs
	git clone --depth 1 git@github.com:webbsalad/storya-recs-backend.git $(TMP_DIR)
	mv $(TMP_DIR)/api/recs api/
	rm -rf $(TMP_DIR)

generate:
	@protoc \
		-I . \
		-I vendor.protogen \
		--validate_out="lang=go:./internal/pb" \
		--go_out=./internal/pb \
		--go-grpc_out=./internal/pb \
		--grpc-gateway_out ./internal/pb \
		--grpc-gateway_opt generate_unbound_methods=true \
		--openapiv2_out=./internal/pb \
		./api/passport/*.proto \
		./api/otp/*.proto \
		./api/content/*.proto \
		./api/recs/*.proto \
	&& swagger mixin --ignore-conflicts -o internal/docs/gateway.swagger.json \
		internal/pb/api/passport/passport_service.swagger.json \
		internal/pb/api/otp/otp_service.swagger.json \
		internal/pb/api/content/content_service.swagger.json \
		internal/pb/api/content/user_content_service.swagger.json \
		internal/pb/api/recs/recs_service.swagger.json \

