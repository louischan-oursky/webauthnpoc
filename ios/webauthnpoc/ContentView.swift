//
//  ContentView.swift
//  webauthnpoc
//
//  Created by louischan on 15/7/2022.
//

import SwiftUI
import AuthenticationServices

class MyState: NSObject, ObservableObject, ASWebAuthenticationPresentationContextProviding {
    var session: ASWebAuthenticationSession?

    func open() {
        let session = ASWebAuthenticationSession.init(url: URL(string: "https://webauthn.com")!, callbackURLScheme: nil) { url, error in
            print("url: \(url)")
            print("error: \(error)")
            self.session = nil
        }
        session.prefersEphemeralWebBrowserSession = true
        session.presentationContextProvider = self
        session.start()
        self.session = session
    }

    func presentationAnchor(for session: ASWebAuthenticationSession) -> ASPresentationAnchor {
        return UIApplication.shared.windows.filter {$0.isKeyWindow}.first!
    }
}

struct ContentView: View {
    @StateObject var state = MyState()

    var body: some View {
        Button("Open ASWebAuthenticationSession") {
            state.open()
        }
    }
}

struct ContentView_Previews: PreviewProvider {
    static var previews: some View {
        ContentView()
    }
}
