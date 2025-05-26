# Java compiler
JC = javac

# Java runtime
JAVA = java

# Source files
SOURCES = Main.java Worker.java

# Class files
CLASSES = Main.class Worker.class

# Default target
all: $(CLASSES)

# Compile .java files to .class files
%.class: %.java
	$(JC) $(SOURCES)

# Run the application
run: $(CLASSES)
	$(JAVA) Main

# Clean generated files
clean:
	rm -f *.class

.PHONY: all run clean
