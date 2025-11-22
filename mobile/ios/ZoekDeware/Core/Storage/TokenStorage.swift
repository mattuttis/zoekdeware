import Foundation
import Security

actor TokenStorage {
    static let shared = TokenStorage()

    private let accessTokenKey = "com.zoekdeware.accessToken"
    private let refreshTokenKey = "com.zoekdeware.refreshToken"

    var accessToken: String? {
        get { read(key: accessTokenKey) }
    }

    var refreshToken: String? {
        get { read(key: refreshTokenKey) }
    }

    func store(accessToken: String, refreshToken: String) {
        save(key: accessTokenKey, value: accessToken)
        save(key: refreshTokenKey, value: refreshToken)
    }

    func clear() {
        delete(key: accessTokenKey)
        delete(key: refreshTokenKey)
    }

    private func save(key: String, value: String) {
        let data = value.data(using: .utf8)!

        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrAccount as String: key,
            kSecValueData as String: data
        ]

        SecItemDelete(query as CFDictionary)
        SecItemAdd(query as CFDictionary, nil)
    }

    private func read(key: String) -> String? {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrAccount as String: key,
            kSecReturnData as String: true,
            kSecMatchLimit as String: kSecMatchLimitOne
        ]

        var result: AnyObject?
        let status = SecItemCopyMatching(query as CFDictionary, &result)

        guard status == errSecSuccess,
              let data = result as? Data,
              let value = String(data: data, encoding: .utf8) else {
            return nil
        }

        return value
    }

    private func delete(key: String) {
        let query: [String: Any] = [
            kSecClass as String: kSecClassGenericPassword,
            kSecAttrAccount as String: key
        ]
        SecItemDelete(query as CFDictionary)
    }
}
