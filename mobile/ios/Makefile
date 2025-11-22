.PHONY: setup generate open clean

# Install XcodeGen if not present and generate project
setup:
	@which xcodegen > /dev/null || brew install xcodegen
	@$(MAKE) generate

# Generate Xcode project from project.yml
generate:
	xcodegen generate

# Open project in Xcode
open:
	open ZoekDeware.xcodeproj

# Clean generated files
clean:
	rm -rf ZoekDeware.xcodeproj
	rm -rf DerivedData

# Build for simulator
build:
	xcodebuild -project ZoekDeware.xcodeproj \
		-scheme ZoekDeware \
		-sdk iphonesimulator \
		-destination 'platform=iOS Simulator,name=iPhone 15' \
		build

# Run tests
test:
	xcodebuild -project ZoekDeware.xcodeproj \
		-scheme ZoekDeware \
		-sdk iphonesimulator \
		-destination 'platform=iOS Simulator,name=iPhone 15' \
		test
