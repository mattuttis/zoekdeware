import Foundation

enum AuthEndpoint: Endpoint {
    case register(email: String, password: String)
    case login(email: String, password: String)
    case refresh(token: String)

    var path: String {
        switch self {
        case .register: return "/auth/register"
        case .login: return "/auth/login"
        case .refresh: return "/auth/refresh"
        }
    }

    var method: HTTPMethod {
        .POST
    }

    var body: Encodable? {
        switch self {
        case .register(let email, let password):
            return ["email": email, "password": password]
        case .login(let email, let password):
            return ["email": email, "password": password]
        case .refresh(let token):
            return ["refresh_token": token]
        }
    }
}

enum ProfileEndpoint: Endpoint {
    case get
    case update(UpdateProfileRequest)

    var path: String { "/profile" }

    var method: HTTPMethod {
        switch self {
        case .get: return .GET
        case .update: return .PUT
        }
    }

    var body: Encodable? {
        switch self {
        case .get: return nil
        case .update(let request): return request
        }
    }
}

enum DiscoverEndpoint: Endpoint {
    case getProfiles(limit: Int)
    case swipe(swipedId: String, direction: SwipeDirection)
    case matches

    var path: String {
        switch self {
        case .getProfiles(let limit): return "/discover?limit=\(limit)"
        case .swipe: return "/swipe"
        case .matches: return "/matches"
        }
    }

    var method: HTTPMethod {
        switch self {
        case .getProfiles, .matches: return .GET
        case .swipe: return .POST
        }
    }

    var body: Encodable? {
        switch self {
        case .swipe(let swipedId, let direction):
            return SwipeRequest(swipedId: swipedId, direction: direction)
        default:
            return nil
        }
    }
}

enum ChatEndpoint: Endpoint {
    case conversations
    case conversation(id: String)
    case sendMessage(conversationId: String, content: String)

    var path: String {
        switch self {
        case .conversations: return "/conversations"
        case .conversation(let id): return "/conversations/\(id)"
        case .sendMessage(let id, _): return "/conversations/\(id)/messages"
        }
    }

    var method: HTTPMethod {
        switch self {
        case .conversations, .conversation: return .GET
        case .sendMessage: return .POST
        }
    }

    var body: Encodable? {
        switch self {
        case .sendMessage(_, let content):
            return ["content": content]
        default:
            return nil
        }
    }
}
