import Foundation

struct Conversation: Codable, Identifiable {
    let id: String
    let participant: ConversationParticipant
    let lastMessage: Message?
    let unreadCount: Int
}

struct ConversationParticipant: Codable {
    let id: String
    let displayName: String
    let photoUrl: String?
}

struct Message: Codable, Identifiable {
    let id: String
    let senderId: String
    let content: String
    let sentAt: Date
    let readAt: Date?

    var isRead: Bool { readAt != nil }
}

struct ConversationsResponse: Codable {
    let conversations: [Conversation]
}

struct ConversationResponse: Codable {
    let conversation: Conversation
    let messages: [Message]
}
