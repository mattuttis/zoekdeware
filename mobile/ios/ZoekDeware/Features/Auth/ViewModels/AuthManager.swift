import Foundation
import SwiftUI

@MainActor
class AuthManager: ObservableObject {
    @Published var isAuthenticated = false
    @Published var isLoading = false
    @Published var error: String?

    init() {
        Task {
            await checkAuthentication()
        }
    }

    func checkAuthentication() async {
        if await TokenStorage.shared.accessToken != nil {
            isAuthenticated = true
        }
    }

    func login(email: String, password: String) async {
        isLoading = true
        error = nil

        do {
            let response: AuthResponse = try await APIClient.shared.request(
                AuthEndpoint.login(email: email, password: password)
            )
            await TokenStorage.shared.store(
                accessToken: response.accessToken,
                refreshToken: response.refreshToken
            )
            isAuthenticated = true
        } catch {
            self.error = error.localizedDescription
        }

        isLoading = false
    }

    func register(email: String, password: String) async {
        isLoading = true
        error = nil

        do {
            let response: AuthResponse = try await APIClient.shared.request(
                AuthEndpoint.register(email: email, password: password)
            )
            await TokenStorage.shared.store(
                accessToken: response.accessToken,
                refreshToken: response.refreshToken
            )
            isAuthenticated = true
        } catch {
            self.error = error.localizedDescription
        }

        isLoading = false
    }

    func logout() async {
        await TokenStorage.shared.clear()
        isAuthenticated = false
    }
}
