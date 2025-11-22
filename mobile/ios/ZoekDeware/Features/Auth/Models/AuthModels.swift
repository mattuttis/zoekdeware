import Foundation

struct AuthResponse: Codable {
    let accessToken: String
    let refreshToken: String
    let expiresAt: Int64
}
