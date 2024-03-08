from cucinassistant.exceptions import CAError
from cucinassistant.database import *


class TestUsers:
    def __init__(self, tester, fake_user):
        self.t = tester
        self.fake_user = fake_user

        # Executes all the subtests IN ORDER
        for name in sorted(dir(self)):
            if name[0] == 'S':
                with self.t.subTest(sub=name[:3]):
                    getattr(self, name)()

    def S00_create_user(self):
        # Tests for create_user
        self.t.assertRaisesRegex(CAError, 'Nome utente non valido \(lunghezza minima', create_user, '', '', '')
        self.t.assertRaisesRegex(CAError, 'Nome utente non valido \(solo lettere', create_user, 'antonio&', '', '')
        self.t.assertRaisesRegex(CAError, 'Password non valida \(lunghezza minima', create_user, 'antonio', '', '')
        self.antonio = create_user('antonio', 'antonio@email.com', 'passwordA')
        self.t.assertRaisesRegex(CAError, 'Email non disponibile', create_user, 'antonello', 'antonio@email.com', 'passwordA')
        self.t.assertRaisesRegex(CAError, 'Nome utente non disponibile', create_user, 'antonio', 'antonello@email.com', 'passwordA')
        self.antonello = create_user('antonello', 'antonello@email.com', 'passwordA')

    def S01_login(self):
        # Tests for login
        self.t.assertEqual(self.antonio, login('antonio', 'passwordA'))
        self.t.assertRaisesRegex(CAError, 'Credenziali non valide', login, 'antonio', 'passwordB')
        self.t.assertRaisesRegex(CAError, 'Credenziali non valide', login, 'antonino', 'passwordB')

    def S02_get_data(self):
        # Tests for get_data
        antonio_data = get_data(self.antonio)
        self.t.assertEqual(antonio_data.uid, self.antonio)
        self.t.assertEqual(antonio_data.username, 'antonio')
        self.t.assertEqual(antonio_data.email, 'antonio@email.com')
        self.t.assertEqual(antonio_data, get_data(0, email='antonio@email.com'))
        self.t.assertRaisesRegex(CAError, 'Utente sconosciuto', get_data, self.fake_user)

    def S03_change_username(self):
        # Tests for change_username
        self.t.assertRaisesRegex(CAError, 'Utente sconosciuto', change_username, self.fake_user, 'nuovo_username')
        change_username(self.antonello, 'antonello')
        self.t.assertEqual(get_data(self.antonello).username, 'antonello')
        self.t.assertRaisesRegex(CAError, 'Nome utente non disponibile', change_username, self.antonello, 'antonio')
        change_username(self.antonello, 'antonellino')
        self.t.assertEqual(get_data(self.antonello).username, 'antonellino')

    def S04_change_email(self):
        # Tests for change_email
        self.t.assertRaisesRegex(CAError, 'Utente sconosciuto', change_email, self.fake_user, 'fake_email')
        change_email(self.antonello, 'antonello@email.com')
        self.t.assertEqual(get_data(self.antonello).email, 'antonello@email.com')
        self.t.assertRaisesRegex(CAError, 'Email non disponibile', change_email, self.antonello, 'antonio@email.com')
        change_email(self.antonello, 'antonellino@email.com')
        self.t.assertEqual(get_data(self.antonello).email, 'antonellino@email.com')

    def S05_change_password(self):
        # Tests for change_password
        self.t.assertRaisesRegex(CAError, 'Utente sconosciuto', change_password, self.fake_user, 'fake_password', 'fake_password2')
        self.t.assertRaisesRegex(CAError, 'Credenziali non valide', change_password, self.antonello, 'vecchia', 'nuova')
        change_password(self.antonello, 'passwordA', 'passwordB')
        self.t.assertTrue(ph.verify(get_data(self.antonello).password, 'passwordB'))

    def S06_generate_token(self):
        # Tests for generate_token
        self.token = generate_token(self.antonello)
        self.t.assertTrue(ph.verify(get_data(self.antonello).token, self.token))
        self.t.assertRaisesRegex(CAError, 'Utente sconosciuto', generate_token, self.fake_user)

    def S07_reset_password(self):
        # Test for reset_password
        self.t.assertRaisesRegex(CAError, 'Utente sconosciuto', reset_password, 'fake@email.com', 'fake_token', 'fake_password')
        self.t.assertRaisesRegex(CAError, 'Errore durante la reimpostazione', reset_password, 'antonellino@email.com', 'token', 'new_password')
        reset_password('antonellino@email.com', self.token, 'passwordA')
        self.t.assertTrue(ph.verify(get_data(self.antonello).password, 'passwordA'))
        self.t.assertIsNone(get_data(self.antonello).token)

    def S08_delete_user(self):
        # Tests for delete_user
        self.t.assertRaisesRegex(CAError, 'Utente sconosciuto', delete_user, self.fake_user, 'token')
        self.t.assertRaisesRegex(CAError, 'Errore durante la cancellazione', delete_user, self.antonello, 'token')
        self.t.assertRaisesRegex(CAError, 'Errore durante la cancellazione', delete_user, self.antonello, self.token)
        self.token = generate_token(self.antonello)
        delete_user(self.antonello, self.token)
        self.t.assertRaisesRegex(CAError, 'Credenziali non valide', login, 'antonello', 'passwordA')

    def S09_get_users_number(self):
        # Tests for get_users_number
        self.t.assertEqual(get_users_number(), 3)

    def S10_get_users_email(self):
        # Tests for get_users_email
        self.t.assertEqual(set(get_users_emails()), {'francesco@email.com', 'giovanna@email.com', 'antonio@email.com'})
