import SwiftUI

struct ProfileView: View {
    @EnvironmentObject var authManager: AuthManager
    @StateObject private var viewModel = ProfileViewModel()

    var body: some View {
        NavigationStack {
            List {
                if let profile = viewModel.profile {
                    Section {
                        HStack {
                            Circle()
                                .fill(.gray.opacity(0.3))
                                .frame(width: 80, height: 80)
                                .overlay {
                                    Image(systemName: "person.fill")
                                        .font(.largeTitle)
                                        .foregroundStyle(.gray)
                                }

                            VStack(alignment: .leading, spacing: 4) {
                                Text(profile.displayName ?? "No name")
                                    .font(.title2.bold())

                                if let age = profile.age {
                                    Text("\(age) years old")
                                        .foregroundStyle(.secondary)
                                }
                            }
                        }
                        .padding(.vertical, 8)
                    }

                    Section("Bio") {
                        Text(profile.bio ?? "No bio yet")
                            .foregroundStyle(profile.bio == nil ? .secondary : .primary)
                    }

                    Section("Interests") {
                        if profile.interests.isEmpty {
                            Text("No interests added")
                                .foregroundStyle(.secondary)
                        } else {
                            FlowLayout(spacing: 8) {
                                ForEach(profile.interests, id: \.self) { interest in
                                    Text(interest)
                                        .padding(.horizontal, 12)
                                        .padding(.vertical, 6)
                                        .background(.pink.opacity(0.1))
                                        .foregroundStyle(.pink)
                                        .clipShape(Capsule())
                                }
                            }
                        }
                    }

                    Section {
                        NavigationLink("Edit Profile") {
                            EditProfileView(viewModel: viewModel)
                        }
                    }

                    Section {
                        Button("Log Out", role: .destructive) {
                            Task {
                                await authManager.logout()
                            }
                        }
                    }
                } else if viewModel.isLoading {
                    ProgressView()
                        .frame(maxWidth: .infinity)
                } else if let error = viewModel.error {
                    Text(error)
                        .foregroundStyle(.red)
                }
            }
            .navigationTitle("Profile")
            .task {
                await viewModel.loadProfile()
            }
            .refreshable {
                await viewModel.loadProfile()
            }
        }
    }
}

struct EditProfileView: View {
    @ObservedObject var viewModel: ProfileViewModel
    @Environment(\.dismiss) var dismiss

    @State private var displayName = ""
    @State private var bio = ""
    @State private var birthDate = Date()
    @State private var gender: Gender?

    var body: some View {
        Form {
            Section("Basic Info") {
                TextField("Display Name", text: $displayName)
                DatePicker("Birthday", selection: $birthDate, displayedComponents: .date)
                Picker("Gender", selection: $gender) {
                    Text("Select").tag(nil as Gender?)
                    ForEach(Gender.allCases, id: \.self) { gender in
                        Text(gender.displayName).tag(gender as Gender?)
                    }
                }
            }

            Section("About") {
                TextEditor(text: $bio)
                    .frame(minHeight: 100)
            }

            Section {
                Button("Save") {
                    Task {
                        await viewModel.updateProfile(
                            displayName: displayName,
                            bio: bio,
                            birthDate: birthDate,
                            gender: gender
                        )
                        dismiss()
                    }
                }
                .disabled(viewModel.isLoading)
            }
        }
        .navigationTitle("Edit Profile")
        .onAppear {
            if let profile = viewModel.profile {
                displayName = profile.displayName ?? ""
                bio = profile.bio ?? ""
                birthDate = profile.birthDate ?? Date()
                gender = profile.gender
            }
        }
    }
}

struct FlowLayout: Layout {
    var spacing: CGFloat = 8

    func sizeThatFits(proposal: ProposedViewSize, subviews: Subviews, cache: inout ()) -> CGSize {
        let result = FlowResult(in: proposal.width ?? 0, subviews: subviews, spacing: spacing)
        return result.size
    }

    func placeSubviews(in bounds: CGRect, proposal: ProposedViewSize, subviews: Subviews, cache: inout ()) {
        let result = FlowResult(in: bounds.width, subviews: subviews, spacing: spacing)
        for (index, subview) in subviews.enumerated() {
            subview.place(at: CGPoint(x: bounds.minX + result.positions[index].x,
                                      y: bounds.minY + result.positions[index].y),
                         proposal: .unspecified)
        }
    }

    struct FlowResult {
        var size: CGSize = .zero
        var positions: [CGPoint] = []

        init(in maxWidth: CGFloat, subviews: Subviews, spacing: CGFloat) {
            var x: CGFloat = 0
            var y: CGFloat = 0
            var rowHeight: CGFloat = 0

            for subview in subviews {
                let size = subview.sizeThatFits(.unspecified)

                if x + size.width > maxWidth, x > 0 {
                    x = 0
                    y += rowHeight + spacing
                    rowHeight = 0
                }

                positions.append(CGPoint(x: x, y: y))
                rowHeight = max(rowHeight, size.height)
                x += size.width + spacing
            }

            self.size = CGSize(width: maxWidth, height: y + rowHeight)
        }
    }
}

#Preview {
    ProfileView()
        .environmentObject(AuthManager())
}
