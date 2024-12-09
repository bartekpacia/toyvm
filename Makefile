EXAMPLES_DIR = examples
NASM_FILES = $(wildcard $(EXAMPLES_DIR)/*.nasm)
BIN_FILES = $(NASM_FILES:$(EXAMPLES_DIR)/%.nasm=$(EXAMPLES_DIR)/%.bin)

all: $(BIN_FILES)

$(EXAMPLES_DIR)/%.bin: $(EXAMPLES_DIR)/%.nasm
	cd $(EXAMPLES_DIR) && nasm -o $(<F:.nasm=.bin) $(<F)

clean:
	rm -f $(BIN_FILES)
