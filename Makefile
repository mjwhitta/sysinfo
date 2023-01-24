-include gomk/main.mk
-include local/Makefile

ifneq ($(unameS),Windows)
spellcheck:
	@codespell -f -L hilight,hilights -S ".git,*.pem"
endif
