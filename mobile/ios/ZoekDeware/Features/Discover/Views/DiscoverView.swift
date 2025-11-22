import SwiftUI

struct DiscoverView: View {
    @StateObject private var viewModel = DiscoverViewModel()

    var body: some View {
        NavigationStack {
            ZStack {
                if viewModel.isLoading && viewModel.profiles.isEmpty {
                    ProgressView("Finding people near you...")
                } else if let error = viewModel.error, viewModel.profiles.isEmpty {
                    VStack(spacing: 16) {
                        Image(systemName: "wifi.slash")
                            .font(.largeTitle)
                            .foregroundStyle(.secondary)
                        Text(error)
                            .foregroundStyle(.secondary)
                        Button("Try Again") {
                            Task { await viewModel.loadProfiles() }
                        }
                        .buttonStyle(.borderedProminent)
                    }
                } else if let profile = viewModel.currentProfile {
                    CardStackView(profile: profile, viewModel: viewModel)
                } else {
                    VStack(spacing: 16) {
                        Image(systemName: "heart.slash")
                            .font(.system(size: 64))
                            .foregroundStyle(.secondary)
                        Text("No more profiles")
                            .font(.title2)
                        Text("Check back later for new people")
                            .foregroundStyle(.secondary)
                        Button("Refresh") {
                            Task { await viewModel.loadProfiles() }
                        }
                        .buttonStyle(.borderedProminent)
                        .tint(.pink)
                    }
                }
            }
            .navigationTitle("Discover")
            .task {
                await viewModel.loadProfiles()
            }
            .sheet(item: $viewModel.newMatch) { match in
                MatchView(match: match) {
                    viewModel.dismissMatch()
                }
            }
        }
    }
}

struct CardStackView: View {
    let profile: DiscoverProfile
    @ObservedObject var viewModel: DiscoverViewModel
    @State private var offset: CGSize = .zero
    @State private var rotation: Double = 0

    var body: some View {
        VStack {
            ZStack {
                ProfileCardView(profile: profile)
                    .offset(offset)
                    .rotationEffect(.degrees(rotation))
                    .gesture(
                        DragGesture()
                            .onChanged { gesture in
                                offset = gesture.translation
                                rotation = Double(gesture.translation.width / 20)
                            }
                            .onEnded { gesture in
                                if abs(gesture.translation.width) > 150 {
                                    let direction: SwipeDirection = gesture.translation.width > 0 ? .like : .pass
                                    withAnimation(.spring()) {
                                        offset = CGSize(
                                            width: gesture.translation.width > 0 ? 500 : -500,
                                            height: 0
                                        )
                                    }
                                    Task {
                                        await viewModel.swipe(direction: direction)
                                        offset = .zero
                                        rotation = 0
                                    }
                                } else {
                                    withAnimation(.spring()) {
                                        offset = .zero
                                        rotation = 0
                                    }
                                }
                            }
                    )
            }

            HStack(spacing: 32) {
                SwipeButton(systemName: "xmark", color: .red) {
                    Task { await viewModel.swipe(direction: .pass) }
                }

                SwipeButton(systemName: "star.fill", color: .blue) {
                    Task { await viewModel.swipe(direction: .superLike) }
                }

                SwipeButton(systemName: "heart.fill", color: .green) {
                    Task { await viewModel.swipe(direction: .like) }
                }
            }
            .padding(.top, 24)
        }
        .padding()
    }
}

struct ProfileCardView: View {
    let profile: DiscoverProfile

    var body: some View {
        VStack(alignment: .leading, spacing: 0) {
            ZStack(alignment: .bottomLeading) {
                Rectangle()
                    .fill(.gray.opacity(0.3))
                    .aspectRatio(0.75, contentMode: .fit)
                    .overlay {
                        if let photoUrl = profile.photos.first {
                            AsyncImage(url: URL(string: photoUrl)) { image in
                                image.resizable().scaledToFill()
                            } placeholder: {
                                ProgressView()
                            }
                        } else {
                            Image(systemName: "person.fill")
                                .font(.system(size: 80))
                                .foregroundStyle(.gray)
                        }
                    }
                    .clipped()

                LinearGradient(
                    colors: [.clear, .black.opacity(0.7)],
                    startPoint: .center,
                    endPoint: .bottom
                )

                VStack(alignment: .leading, spacing: 4) {
                    HStack(alignment: .firstTextBaseline) {
                        Text(profile.displayName)
                            .font(.title.bold())
                        Text("\(profile.age)")
                            .font(.title2)
                    }

                    if let distance = profile.distanceKm {
                        Label("\(Int(distance)) km away", systemImage: "location.fill")
                            .font(.subheadline)
                    }
                }
                .foregroundStyle(.white)
                .padding()
            }

            if let bio = profile.bio {
                Text(bio)
                    .padding()
                    .lineLimit(3)
            }
        }
        .background(.background)
        .clipShape(RoundedRectangle(cornerRadius: 16))
        .shadow(radius: 8)
    }
}

struct SwipeButton: View {
    let systemName: String
    let color: Color
    let action: () -> Void

    var body: some View {
        Button(action: action) {
            Image(systemName: systemName)
                .font(.title)
                .foregroundStyle(color)
                .frame(width: 60, height: 60)
                .background(color.opacity(0.1))
                .clipShape(Circle())
        }
    }
}

struct MatchView: View {
    let match: Match
    let onDismiss: () -> Void

    var body: some View {
        VStack(spacing: 24) {
            Spacer()

            Text("It's a Match!")
                .font(.largeTitle.bold())
                .foregroundStyle(.pink)

            Text("You and \(match.displayName) liked each other")
                .foregroundStyle(.secondary)

            Circle()
                .fill(.gray.opacity(0.3))
                .frame(width: 120, height: 120)
                .overlay {
                    Image(systemName: "person.fill")
                        .font(.system(size: 48))
                        .foregroundStyle(.gray)
                }

            Spacer()

            Button("Send a Message") {
                onDismiss()
            }
            .buttonStyle(.borderedProminent)
            .tint(.pink)

            Button("Keep Swiping") {
                onDismiss()
            }
            .foregroundStyle(.secondary)

            Spacer()
        }
        .padding()
    }
}

#Preview {
    DiscoverView()
}
