import Foundation
import SwiftUI

@MainActor
class ProfileViewModel: ObservableObject {
    @Published var profile: Profile?
    @Published var isLoading = false
    @Published var error: String?

    func loadProfile() async {
        isLoading = true
        error = nil

        do {
            profile = try await APIClient.shared.request(ProfileEndpoint.get)
        } catch {
            self.error = error.localizedDescription
        }

        isLoading = false
    }

    func updateProfile(displayName: String, bio: String, birthDate: Date?, gender: Gender?) async {
        isLoading = true
        error = nil

        let request = UpdateProfileRequest(
            displayName: displayName,
            bio: bio,
            birthDate: birthDate,
            gender: gender?.rawValue,
            interests: nil
        )

        do {
            profile = try await APIClient.shared.request(ProfileEndpoint.update(request))
        } catch {
            self.error = error.localizedDescription
        }

        isLoading = false
    }
}
