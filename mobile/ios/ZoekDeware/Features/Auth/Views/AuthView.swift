import SwiftUI

struct AuthView: View {
    @State private var isLogin = true

    var body: some View {
        NavigationStack {
            VStack(spacing: 32) {
                Spacer()

                VStack(spacing: 8) {
                    Image(systemName: "heart.fill")
                        .font(.system(size: 64))
                        .foregroundStyle(.pink)

                    Text("ZoekDeware")
                        .font(.largeTitle.bold())
                }

                if isLogin {
                    LoginView()
                } else {
                    RegisterView()
                }

                Button {
                    withAnimation {
                        isLogin.toggle()
                    }
                } label: {
                    Text(isLogin ? "Don't have an account? Sign up" : "Already have an account? Log in")
                        .foregroundStyle(.secondary)
                }

                Spacer()
            }
            .padding()
        }
    }
}

struct LoginView: View {
    @EnvironmentObject var authManager: AuthManager
    @State private var email = ""
    @State private var password = ""

    var body: some View {
        VStack(spacing: 16) {
            TextField("Email", text: $email)
                .textFieldStyle(.roundedBorder)
                .textContentType(.emailAddress)
                .autocapitalization(.none)
                .keyboardType(.emailAddress)

            SecureField("Password", text: $password)
                .textFieldStyle(.roundedBorder)
                .textContentType(.password)

            if let error = authManager.error {
                Text(error)
                    .foregroundStyle(.red)
                    .font(.caption)
            }

            Button {
                Task {
                    await authManager.login(email: email, password: password)
                }
            } label: {
                if authManager.isLoading {
                    ProgressView()
                        .frame(maxWidth: .infinity)
                } else {
                    Text("Log In")
                        .frame(maxWidth: .infinity)
                }
            }
            .buttonStyle(.borderedProminent)
            .tint(.pink)
            .disabled(email.isEmpty || password.isEmpty || authManager.isLoading)
        }
    }
}

struct RegisterView: View {
    @EnvironmentObject var authManager: AuthManager
    @State private var email = ""
    @State private var password = ""
    @State private var confirmPassword = ""

    var passwordsMatch: Bool {
        password == confirmPassword && !password.isEmpty
    }

    var body: some View {
        VStack(spacing: 16) {
            TextField("Email", text: $email)
                .textFieldStyle(.roundedBorder)
                .textContentType(.emailAddress)
                .autocapitalization(.none)
                .keyboardType(.emailAddress)

            SecureField("Password", text: $password)
                .textFieldStyle(.roundedBorder)
                .textContentType(.newPassword)

            SecureField("Confirm Password", text: $confirmPassword)
                .textFieldStyle(.roundedBorder)
                .textContentType(.newPassword)

            if let error = authManager.error {
                Text(error)
                    .foregroundStyle(.red)
                    .font(.caption)
            }

            Button {
                Task {
                    await authManager.register(email: email, password: password)
                }
            } label: {
                if authManager.isLoading {
                    ProgressView()
                        .frame(maxWidth: .infinity)
                } else {
                    Text("Sign Up")
                        .frame(maxWidth: .infinity)
                }
            }
            .buttonStyle(.borderedProminent)
            .tint(.pink)
            .disabled(!passwordsMatch || email.isEmpty || authManager.isLoading)
        }
    }
}

#Preview {
    AuthView()
        .environmentObject(AuthManager())
}
