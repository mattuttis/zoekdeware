import Foundation

enum Config {
    enum Environment {
        case development
        case staging
        case production
    }

    static let current: Environment = {
        #if DEBUG
        return .development
        #else
        return .production
        #endif
    }()

    // For physical iPhone testing, replace localhost with your Mac's IP address
    // Find it with: ipconfig getifaddr en0
    private static let devHost = "192.168.68.62" // Your Mac's IP for iPhone testing
    private static let devPort = "80" // K8s gateway port

    static var apiBaseURL: URL {
        switch current {
        case .development:
            return URL(string: "http://\(devHost):\(devPort)/api/v1/")!
        case .staging:
            return URL(string: "https://staging-api.zoekdeware.app/api/v1/")!
        case .production:
            return URL(string: "https://api.zoekdeware.app/api/v1/")!
        }
    }

    static var webSocketURL: URL {
        switch current {
        case .development:
            return URL(string: "ws://\(devHost):\(devPort)/api/v1/ws")!
        case .staging:
            return URL(string: "wss://staging-api.zoekdeware.app/api/v1/ws")!
        case .production:
            return URL(string: "wss://api.zoekdeware.app/api/v1/ws")!
        }
    }
}
