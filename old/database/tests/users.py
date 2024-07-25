    def S03_change_username(self):
        # Tests for change_username
        self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.change_username, self.fake_user, 'nuovo_username')
        db.change_username(self.antonello, 'antonello')
        self.t.assertEqual(db.get_data(self.antonello).username, 'antonello')
        self.t.assertRaisesRegex(CAError, 'Nome utente non disponibile', db.change_username, self.antonello, 'antonio')
        db.change_username(self.antonello, 'antonellino')
        self.t.assertEqual(db.get_data(self.antonello).username, 'antonellino')

    def S04_change_email(self):
        # Tests for change_email
        self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.change_email, self.fake_user, 'fake_email')
        db.change_email(self.antonello, 'antonello@email.com')
        self.t.assertEqual(db.get_data(self.antonello).email, 'antonello@email.com')
        self.t.assertRaisesRegex(CAError, 'Email non disponibile', db.change_email, self.antonello, 'antonio@email.com')
        db.change_email(self.antonello, 'antonellino@email.com')
        self.t.assertEqual(db.get_data(self.antonello).email, 'antonellino@email.com')

    def S05_change_password(self):
        # Tests for change_password
        self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.change_password, self.fake_user, 'fake_password', 'fake_password2')
        self.t.assertRaisesRegex(CAError, 'Credenziali non valide', db.change_password, self.antonello, 'vecchia', 'nuova')
        db.change_password(self.antonello, 'passwordA', 'passwordB')
        self.t.assertTrue(db.ph.verify(db.get_data(self.antonello).password, 'passwordB'))

    def S07_reset_password(self):
        # Test for reset_password
        self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.reset_password, 'fake@email.com', 'fake_token', 'fake_password')
        self.t.assertRaisesRegex(CAError, 'Errore durante la reimpostazione', db.reset_password, 'antonellino@email.com', 'token', 'new_password')
        db.reset_password('antonellino@email.com', self.token, 'passwordA')
        self.t.assertTrue(db.ph.verify(db.get_data(self.antonello).password, 'passwordA'))
        self.t.assertIsNone(db.get_data(self.antonello).token)

    def S08_delete_user(self):
        # Tests for delete_user
        self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.delete_user, self.fake_user, 'token')
        self.t.assertRaisesRegex(CAError, 'Errore durante la cancellazione', db.delete_user, self.antonello, 'token')
        self.t.assertRaisesRegex(CAError, 'Errore durante la cancellazione', db.delete_user, self.antonello, self.token)
        self.token = db.generate_token(self.antonello)
        db.append_storage(self.antonello, [['storage', '', '']])
        db.append_list(self.antonello, 'shopping', 'lists')
        db.append_list(self.antonello, 'ideas', 'lists')
        db.create_menu(self.antonello)
        db.duplicate_menu(self.antonello, 1)
        db.duplicate_menu(self.antonello, 2)
        db.delete_user(self.antonello, self.token)
        self.t.assertRaisesRegex(CAError, 'Credenziali non valide', db.login, 'antonello', 'passwordA')

    def S10_get_users_email(self):
        # Tests for get_users_email
        self.t.assertEqual(set(db.get_users_emails()), {'francesco@email.com', 'giovanna@email.com', 'antonio@email.com'})
