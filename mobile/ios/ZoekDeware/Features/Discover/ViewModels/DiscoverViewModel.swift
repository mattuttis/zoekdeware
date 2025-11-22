import Foundation
import SwiftUI

@MainActor
class DiscoverViewModel: ObservableObject {
    @Published var profiles: [DiscoverProfile] = []
    @Published var currentIndex = 0
    @Published var isLoading = false
    @Published var error: String?
    @Published var newMatch: Match?

    var currentProfile: DiscoverProfile? {
        guard currentIndex < profiles.count else { return nil }
        return profiles[currentIndex]
    }

    func loadProfiles() async {
        guard !isLoading else { return }

        isLoading = true
        error = nil

        do {
            let response: DiscoverResponse = try await APIClient.shared.request(
                DiscoverEndpoint.getProfiles(limit: 10)
            )
            profiles = response.profiles
            currentIndex = 0
        } catch {
            self.error = error.localizedDescription
        }

        isLoading = false
    }

    func swipe(direction: SwipeDirection) async {
        guard let profile = currentProfile else { return }

        do {
            let response: SwipeResponse = try await APIClient.shared.request(
                DiscoverEndpoint.swipe(swipedId: profile.id, direction: direction)
            )

            if let match = response.match {
                newMatch = match
            }

            withAnimation {
                currentIndex += 1
            }

            if currentIndex >= profiles.count - 2 {
                await loadProfiles()
            }
        } catch {
            self.error = error.localizedDescription
        }
    }

    func dismissMatch() {
        newMatch = nil
    }
}

struct DiscoverResponse: Codable {
    let profiles: [DiscoverProfile]
}
