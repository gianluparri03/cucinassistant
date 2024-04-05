from cucinassistant.exceptions import CAError, CACritical
import cucinassistant.database as db


class TestUsers:
    def __init__(self, tester, fake_user):
        self.t = tester
        self.fake_user = fake_user

        # Executes all the subtests IN ORDER
        for name in sorted(dir(self)):
            if name[0] == 'S':
                with self.t.subTest(cat='Users', sub=name[:3]):
                    getattr(self, name)()

    def S00_create_user(self):
        # Tests for create_user
        self.t.assertRaisesRegex(CAError, r'Nome utente non valido \(lunghezza minima', db.create_user, '', '', '')
        self.t.assertRaisesRegex(CAError, r'Nome utente non valido \(solo lettere', db.create_user, 'antonio&', '', '')
        self.t.assertRaisesRegex(CAError, r'Password non valida \(lunghezza minima', db.create_user, 'antonio', '', '')
        self.antonio = db.create_user('antonio', 'antonio@email.com', 'passwordA')
        self.t.assertRaisesRegex(CAError, 'Email non disponibile', db.create_user, 'antonello', 'antonio@email.com', 'passwordA')
        self.t.assertRaisesRegex(CAError, 'Nome utente non disponibile', db.create_user, 'antonio', 'antonello@email.com', 'passwordA')
        self.antonello = db.create_user('antonello', 'antonello@email.com', 'passwordA')

    def S01_login(self):
        # Tests for login
        self.t.assertEqual(self.antonio, db.login('antonio', 'passwordA'))
        self.t.assertRaisesRegex(CAError, 'Credenziali non valide', db.login, 'antonio', 'passwordB')
        self.t.assertRaisesRegex(CAError, 'Credenziali non valide', db.login, 'antonino', 'passwordB')

    def S02_get_data(self):
        # Tests for get_data
        antonio_data = db.get_data(self.antonio)
        self.t.assertEqual(antonio_data.uid, self.antonio)
        self.t.assertEqual(antonio_data.username, 'antonio')
        self.t.assertEqual(antonio_data.email, 'antonio@email.com')
        self.t.assertEqual(antonio_data, db.get_data(0, email='antonio@email.com'))
        self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.get_data, self.fake_user)

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

    def S06_generate_token(self):
        # Tests for generate_token
        self.token = db.generate_token(self.antonello)
        self.t.assertTrue(db.ph.verify(db.get_data(self.antonello).token, self.token))
        self.t.assertRaisesRegex(CACritical, 'Utente sconosciuto', db.generate_token, self.fake_user)

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

    def S09_get_users_number(self):
        # Tests for get_users_number
        self.t.assertEqual(db.get_users_number(), 3)

    def S10_get_users_email(self):
        # Tests for get_users_email
        self.t.assertEqual(set(db.get_users_emails()), {'francesco@email.com', 'giovanna@email.com', 'antonio@email.com'})
