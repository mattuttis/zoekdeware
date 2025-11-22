import Foundation

@MainActor
class ChatListViewModel: ObservableObject {
    @Published var conversations: [Conversation] = []
    @Published var isLoading = false
    @Published var error: String?

    func loadConversations() async {
        isLoading = true
        error = nil

        do {
            let response: ConversationsResponse = try await APIClient.shared.request(
                ChatEndpoint.conversations
            )
            conversations = response.conversations
        } catch {
            self.error = error.localizedDescription
        }

        isLoading = false
    }
}

@MainActor
class ChatDetailViewModel: ObservableObject {
    let conversationId: String

    @Published var conversation: Conversation?
    @Published var messages: [Message] = []
    @Published var isLoading = false
    @Published var isSending = false
    @Published var error: String?

    init(conversationId: String) {
        self.conversationId = conversationId
    }

    func loadConversation() async {
        isLoading = true
        error = nil

        do {
            let response: ConversationResponse = try await APIClient.shared.request(
                ChatEndpoint.conversation(id: conversationId)
            )
            conversation = response.conversation
            messages = response.messages
        } catch {
            self.error = error.localizedDescription
        }

        isLoading = false
    }

    func sendMessage(_ content: String) async {
        guard !content.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty else { return }

        isSending = true

        do {
            let message: Message = try await APIClient.shared.request(
                ChatEndpoint.sendMessage(conversationId: conversationId, content: content)
            )
            messages.append(message)
        } catch {
            self.error = error.localizedDescription
        }

        isSending = false
    }
}
