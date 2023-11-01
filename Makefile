-include gomk/main.mk
-include local/Makefile

ifneq ($(unameS),Windows)
spellcheck:
	@codespell -f -L hilight,hilighter,hilights -S ".git,*.pem"
endif
