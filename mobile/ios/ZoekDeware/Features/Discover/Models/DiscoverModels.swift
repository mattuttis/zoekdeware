import Foundation

struct DiscoverProfile: Codable, Identifiable {
    let id: String
    let displayName: String
    let age: Int
    let bio: String?
    let photos: [String]
    let distanceKm: Double?
}

enum SwipeDirection: String, Codable {
    case like
    case pass
    case superLike = "super_like"
}

struct SwipeRequest: Codable {
    let swipedId: String
    let direction: SwipeDirection
}

struct SwipeResponse: Codable {
    let match: Match?
}

struct Match: Codable, Identifiable {
    let id: String
    let memberId: String
    let displayName: String
    let photoUrl: String?
    let matchedAt: Date
}
