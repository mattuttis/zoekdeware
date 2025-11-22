import SwiftUI

struct ChatListView: View {
    @StateObject private var viewModel = ChatListViewModel()

    var body: some View {
        NavigationStack {
            Group {
                if viewModel.isLoading && viewModel.conversations.isEmpty {
                    ProgressView()
                } else if viewModel.conversations.isEmpty {
                    VStack(spacing: 16) {
                        Image(systemName: "message.fill")
                            .font(.system(size: 48))
                            .foregroundStyle(.secondary)
                        Text("No Conversations")
                            .font(.title2.bold())
                        Text("Match with someone to start chatting")
                            .font(.subheadline)
                            .foregroundStyle(.secondary)
                    }
                } else {
                    List(viewModel.conversations) { conversation in
                        NavigationLink {
                            ChatDetailView(conversationId: conversation.id)
                        } label: {
                            ConversationRow(conversation: conversation)
                        }
                    }
                }
            }
            .navigationTitle("Messages")
            .task {
                await viewModel.loadConversations()
            }
            .refreshable {
                await viewModel.loadConversations()
            }
        }
    }
}

struct ConversationRow: View {
    let conversation: Conversation

    var body: some View {
        HStack(spacing: 12) {
            Circle()
                .fill(.gray.opacity(0.3))
                .frame(width: 56, height: 56)
                .overlay {
                    if let photoUrl = conversation.participant.photoUrl {
                        AsyncImage(url: URL(string: photoUrl)) { image in
                            image.resizable().scaledToFill()
                        } placeholder: {
                            Image(systemName: "person.fill")
                                .foregroundStyle(.gray)
                        }
                        .clipShape(Circle())
                    } else {
                        Image(systemName: "person.fill")
                            .foregroundStyle(.gray)
                    }
                }

            VStack(alignment: .leading, spacing: 4) {
                HStack {
                    Text(conversation.participant.displayName)
                        .font(.headline)

                    Spacer()

                    if let lastMessage = conversation.lastMessage {
                        Text(lastMessage.sentAt, style: .relative)
                            .font(.caption)
                            .foregroundStyle(.secondary)
                    }
                }

                HStack {
                    if let lastMessage = conversation.lastMessage {
                        Text(lastMessage.content)
                            .font(.subheadline)
                            .foregroundStyle(.secondary)
                            .lineLimit(1)
                    } else {
                        Text("Start a conversation")
                            .font(.subheadline)
                            .foregroundStyle(.secondary)
                            .italic()
                    }

                    Spacer()

                    if conversation.unreadCount > 0 {
                        Text("\(conversation.unreadCount)")
                            .font(.caption.bold())
                            .foregroundStyle(.white)
                            .padding(.horizontal, 8)
                            .padding(.vertical, 2)
                            .background(.pink)
                            .clipShape(Capsule())
                    }
                }
            }
        }
        .padding(.vertical, 4)
    }
}

#Preview {
    ChatListView()
}
