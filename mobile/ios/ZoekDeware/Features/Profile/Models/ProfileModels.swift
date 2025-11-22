import Foundation

struct Profile: Codable, Identifiable {
    let id: String
    let email: String
    var displayName: String?
    var bio: String?
    var birthDate: Date?
    var gender: Gender?
    var interests: [String]
    var photos: [String]

    var age: Int? {
        guard let birthDate else { return nil }
        let calendar = Calendar.current
        let now = Date()
        let ageComponents = calendar.dateComponents([.year], from: birthDate, to: now)
        return ageComponents.year
    }
}

enum Gender: String, Codable, CaseIterable {
    case male
    case female
    case other

    var displayName: String {
        switch self {
        case .male: return "Male"
        case .female: return "Female"
        case .other: return "Other"
        }
    }
}

struct UpdateProfileRequest: Codable {
    var displayName: String?
    var bio: String?
    var birthDate: Date?
    var gender: String?
    var interests: [String]?
}
